package connection

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"time"

	"github.com/lorenzoc25/bittorrent-go/client"
	"github.com/lorenzoc25/bittorrent-go/message"
	"github.com/lorenzoc25/bittorrent-go/peer"
)

const MaxBlockSize = 16384

const MaxBacklog = 5

// torrent holds data required to download a torrent from a list of peers
type Torrent struct {
	Peers       []peer.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	data  []byte
}

type pieceProgress struct {
	index      int
	client     *client.Client
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}

func (t *Torrent) getPieceBound(index int) (begin, end int) {
	begin = index * t.PieceLength
	end = begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return
}
func (t *Torrent) getPieceSize(index int) (size int) {
	begin, end := t.getPieceBound(index)
	return end - begin
}

func (state *pieceProgress) readMessgae() error {
	msg, err := state.client.Read()
	if err != nil {
		return err
	}
	if msg == nil {
		return nil
	}
	switch msg.ID {
	case message.MsgChoke:
		state.client.Choked = true
	case message.MsgUnchoke:
		state.client.Choked = false
	case message.MsgHave:
		index, err := message.ParseHave(msg)
		if err != nil {
			return err
		}
		state.client.Bitfield.SetPiece(index)
	case message.MsgPiece:
		n, err := message.ParsePiece(state.index, state.buf, msg)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil
}

func (t *Torrent) Download() ([]byte, error) {
	// init queues for workers to retrive work and send results
	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	resultQueue := make(chan *pieceResult)
	for index, hash := range t.PieceHashes {
		length := t.getPieceSize(index)
		workQueue <- &pieceWork{index, hash, length}
	}
	// start workers
	for _, peer := range t.Peers {
		go t.startDownloadWorker(peer, workQueue, resultQueue)
	}
	// collect all the results
	buf := make([]byte, t.Length)
	donePieces := 0
	for donePieces < len(t.PieceHashes) {
		result := <-resultQueue
		begin, end := t.getPieceBound(result.index)
		copy(buf[begin:end], result.data)
		donePieces++
	}
	close(workQueue)
	return buf, nil
}

// start a worker to download a piece
func attemptDownloadPiece(c *client.Client, pw *pieceWork) ([]byte, error) {
	state := pieceProgress{
		index:  pw.index,
		client: c,
		buf:    make([]byte, pw.length),
	}
	c.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.Conn.SetDeadline(time.Time{})

	for state.downloaded < pw.length {
		// if unchocked, send requsts until we have enough unfullfilled requests
		if !state.client.Choked {
			for state.backlog < MaxBacklog && state.requested < pw.length {
				blockSize := MaxBlockSize
				if pw.length-state.requested < MaxBlockSize {
					blockSize = pw.length - state.requested
				}
				err := c.SendRequest(pw.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize
			}
		}
		err := state.readMessgae()
		if err != nil {
			return nil, err
		}
	}
	return state.buf, nil
}

func checkIntegrity(pw *pieceWork, buf []byte) error {
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], pw.hash[:]) {
		return fmt.Errorf("piece %d hash mismatch", pw.index)
	}
	return nil
}

// start a worker to download a piece
func (t *Torrent) startDownloadWorker(peer peer.Peer, workQueue chan *pieceWork, resultQueue chan *pieceResult) {
	c, err := client.New(peer, t.PeerID, t.InfoHash)
	if err != nil {
		log.Printf("Counld not handshake with peer %s: %s", peer.String(), peer.IP)
		return
	}
	log.Printf("Compelted handshake with peer %s", peer.IP)
	defer c.Conn.Close()

	// choke the connectino and show interest
	c.SendChoke()
	c.SendInterested()

	for piece := range workQueue {
		if !c.Bitfield.HasPiece(piece.index) {
			// put the piece back to workQueue
			workQueue <- piece
			continue
		}
		// start downloading the piece
		buf, err := attemptDownloadPiece(c, piece)
		if err != nil {
			// attemp unsuccessful, put the piece back to workQueue
			workQueue <- piece
			return
		}

		err = checkIntegrity(piece, buf)
		if err != nil {
			// piece failed integrity check, put the piece back to workQueue
			log.Printf("Piece %d failed integrity check", piece.index)
			workQueue <- piece
			return
		}
		// after the piece is downloaded, update other peers that we have the piece
		c.SendHave(piece.index)
		resultQueue <- &pieceResult{piece.index, buf}
	}
}

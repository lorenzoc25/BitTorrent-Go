package connection

import (
	"log"

	"github.com/lorenzoc25/bittorrent-go/client"
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
}

// start a worker to download a piece
func (t *Torrent) startDownloadWorker(peer peer.Peer, workQueue chan *pieceWork, resultQueue chan *pieceResult) {
	c, err := client.New(peer, t.PeerID, t.InfoHash)
	if err != nil {
		log.Printf("Counld not handshake with peer %s: %s", peer.String(), peer.IP)
		return
	}
	defer c.Conn.Close()
}

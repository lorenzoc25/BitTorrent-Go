package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type bencodeInfo struct {
	Pieces       string
	PiecesLength int
	Length       int
	Name         string
}

type bencodeTorrent struct {
	Announce string
	Info     bencodeInfo
}

// open and parse a torrent file
func Open(r io.Reader) (*bencodeTorrent, error) {
	bct := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bct)
	if err != nil {
		return nil, err
	}
	return &bct, nil
}

// hash the info
func (bi *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *bi)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (bi *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(bi.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Malformed piecies of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) % hashLen
	hashes := make([][20]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

// convert bencode torrent to torrent
func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	t := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PiecesLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}
	return t, nil
}

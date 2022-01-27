package handshake

import "io"

type Handshake struct {
	Pstr     string // protocal identifier
	InfoHash [20]byte
	PeerId   [20]byte
}

// serialize the handshake info in a buffer
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	curr_idx := 1
	curr_idx += copy(buf[curr_idx:], h.Pstr)
	curr_idx += copy(buf[curr_idx:], make([]byte, 8))
	curr_idx += copy(buf[curr_idx:], h.InfoHash[:])
	curr_idx += copy(buf[curr_idx:], h.PeerId[:])
	return buf
}

// parses a handshake from return, perform a backward serialize
func Read(r io.Reader) (*Handshake, error) {

}

func New(infoHash, peerId [20]byte) *Handshake {
	return &Handshake{
		InfoHash: infoHash,
		PeerId:   peerId,
	}
}

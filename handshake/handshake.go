package handshake

import (
	"fmt"
	"io"
)

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
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	pstrLen := int(lengthBuf[0])
	if pstrLen == 0 {
		return nil, fmt.Errorf("pstrlen cannot be 0")
	}
	handshakebuf := make([]byte, 48+pstrLen)
	_, err = io.ReadFull(r, handshakebuf)
	if err != nil {
		return nil, err
	}
	var infoHash, peerId [20]byte
	copy(infoHash[:], handshakebuf[pstrLen+8:pstrLen+8+20])
	copy(peerId[:], handshakebuf[pstrLen+8+20:])
	h := Handshake{
		Pstr:     string(handshakebuf[:pstrLen]),
		InfoHash: infoHash,
		PeerId:   peerId,
	}
	return &h, nil
}

func New(infoHash, peerId [20]byte) *Handshake {
	return &Handshake{
		InfoHash: infoHash,
		PeerId:   peerId,
	}
}

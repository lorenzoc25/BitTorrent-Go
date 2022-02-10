package peer

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

// Parses peer ip addresses and ports from a buffer
func Unmarshal(peersBin []byte) ([]Peer, error) {
	const peerSize = 6
	log.Printf("Unmarshal: %d bytes, peerSize is %d", len(peersBin), peerSize)
	numPeers := len(peersBin) / peerSize
	if len(peersBin)%numPeers != 0 {
		err := fmt.Errorf("RECEIVED MALFORMED PEERS")
		return nil, err
	}
	peers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4 : offset+6])
	}
	return peers, nil
}

// stringer method
func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

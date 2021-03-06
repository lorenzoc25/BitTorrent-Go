package client

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/lorenzoc25/bittorrent-go/bitfields"
	"github.com/lorenzoc25/bittorrent-go/handshake"
	"github.com/lorenzoc25/bittorrent-go/message"
	"github.com/lorenzoc25/bittorrent-go/peer"
)

type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield bitfields.Bitfield
	peer     peer.Peer
	infoHash [20]byte
	peerID   [20]byte
}

func compelteHandshake(conn net.Conn, infohash, peerId [20]byte) (*handshake.Handshake, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})
	req := handshake.New(infohash, peerId)
	_, err := conn.Write(req.Serialize())
	if err != nil {
		return nil, err
	}
	res, err := handshake.Read(conn)
	if err != nil {
		return nil, err
	}
	if res.InfoHash != infohash {
		return nil, fmt.Errorf("handshake info hash mismatch")
	}
	return res, nil
}

// receive bitfield from peer
func receiveBitfield(conn net.Conn) (bitfields.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})
	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}
	if msg.ID != message.MsgBitfield {
		err := fmt.Errorf("expected bitfield but got ID %v", msg.ID)
		return nil, err
	}
	return msg.Payload, nil
}

// New connects with a peer, complete a handshake and receive a handshake
func New(peer peer.Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}
	_, err = compelteHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}
	bf, err := receiveBitfield(conn)
	if err != nil {
		return nil, err
	}
	return &Client{
		Conn:     conn,
		Choked:   true,
		Bitfield: bf,
		peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}

// reads and consume a message from the connection
func (c *Client) Read() (*message.Message, error) {
	msg, err := message.Read(c.Conn)
	return msg, err
}

// send request message to peer
func (c *Client) SendRequest(index, begin, length int) error {
	req := message.NewRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}

// make a message to connectino to choke the connection
func (c *Client) SendChoke() error {
	msg := message.Message{ID: message.MsgChoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// make a message to connectino to unchoke the connection
func (c *Client) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// similar to above but for interested
func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// similar to above but for not interested
func (c *Client) SendNotInterested() error {
	msg := message.Message{ID: message.MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// sendHave message to peer
func (c *Client) SendHave(index int) error {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	msg := message.Message{ID: message.MsgHave, Payload: payload}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

package message

// send and receive messages once completed the handshake
import (
	"encoding/binary"
	"io"
)

type messageID uint8

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

// since each message is consist of an id and a payload,
// msg = length + id + payload
type Message struct {
	ID      messageID
	Payload []byte
}

// Create a REQUEST message
func NewRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload, uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:], uint32(length))
	return &Message{
		ID:      MsgRequest,
		Payload: payload,
	}
}

// Create a HAVE message
func NewHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{
		ID:      MsgHave,
		Payload: payload,
	}
}

// serialize the message into a byte array
func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := len(m.Payload) + 1 // +1 for the message id
	buf := make([]byte, length+4)
	binary.BigEndian.PutUint32(buf, uint32(length))
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

// read a message from stream
func Read(r io.Reader) (*Message, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(buf)
	if length == 0 {
		return nil, nil
	}
	if length > 0 {
		buf = make([]byte, length)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
	}
	return &Message{
		ID:      messageID(buf[0]),
		Payload: buf[1:],
	}, nil
}

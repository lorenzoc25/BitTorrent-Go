package bitfields

// bitfield represent the pieces that a peer has
type Bitfield []byte

// HasPiece tells if a bitfield had a piece with given index
func (bf Bitfield) HasPiece(index int) bool {
	return bf[index/8]&(1<<uint(index%8)) != 0
}

// SetPiece sets a piece with given index to 1
func (bf Bitfield) SetPiece(index int) {
	bf[index/8] |= 1 << uint(index%8)
}

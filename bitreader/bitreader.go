package bitreader

type BitReader struct {
	bytes []byte
	curr  int
	nBits uint8
}

func New(bytes []byte) *BitReader {
	return &BitReader{bytes: bytes}
}

// Returns true if popped bit is 1
// false if 0
func (b *BitReader) Pop() bool {
	currByte := b.bytes[b.curr]
	val := (currByte>>b.nBits)&1 == 1

	b.nBits++
	if b.nBits == 8 {
		b.nBits = 0
		b.curr++
	}
	return val
}

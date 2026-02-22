package bitreader

type BitReader struct {
	bytes []byte
	curr  int
	nBits int
}

func New(bytes []byte) *BitReader {
	return &BitReader{bytes: bytes, nBits: 7}
}

// Returns true if popped bit is 1
// false if 0
func (b *BitReader) Pop() bool {
	currByte := b.bytes[b.curr]
	val := (currByte>>b.nBits)&1 == 1

	b.nBits--
	if b.nBits < 0 {
		b.nBits = 7
		b.curr++
	}
	return val
}

func (b *BitReader) HasData() bool {
	return b.curr < len(b.bytes)
}


package bitwriter

type BitWriter struct {
	bytes []byte
	curr  byte
	nBits uint8
}

func New() *BitWriter {
	return &BitWriter{}
}

func (b *BitWriter) WriteUInt(data uint64, size uint8) {
	if size == 0 {
		return
	}

	for i := int(size - 1); i >= 0; i-- {
		bit := (data >> i) & 1
		b.curr = (b.curr << 1) | byte(bit)
		b.nBits++

		if b.nBits == 8 {
			b.bytes = append(b.bytes, b.curr)
			b.curr = 0
			b.nBits = 0
		}
	}
}

func (b *BitWriter) Bytes() []byte {
	if b.nBits > 0 {
		b.curr <<= (8 - b.nBits)
		b.bytes = append(b.bytes, b.curr)
	}
	return b.bytes
}

func (b *BitWriter) TotalBits() int {
	byteCount := len(b.bytes)
	return (byteCount * 8) + int(b.nBits)
}

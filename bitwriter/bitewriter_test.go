package bitwriter

import (
	"bytes"
	"testing"
)



func TestWriteSingleByte(t *testing.T) {
	bw := New()
	bw.WriteUInt(0b10101010, 8)

	result := bw.Bytes()
	expected := []byte{0b10101010}

	if !bytes.Equal(result, expected) {
		t.Fatalf("expected %08b, got %08b", expected, result)
	}
}

func TestWriteAcrossByteBoundary(t *testing.T) {
	bw := New()
	bw.WriteUInt(0b101, 3)
	bw.WriteUInt(0b11100, 5)


	result := bw.Bytes()
	expected := []byte{0b10111100}

	if !bytes.Equal(result, expected) {

		t.Fatalf("expected %08b, got %08b", expected, result)
	}
}

func TestMultipleBytes(t *testing.T) {

	bw := New()
	bw.WriteUInt(0xABCD, 16)


	result := bw.Bytes()
	expected := []byte{0xAB, 0xCD}


	if !bytes.Equal(result, expected) {
		t.Fatalf("expected % X, got % X", expected, result)

	}
}

func TestPaddingOnFlush(t *testing.T) {
	bw := New()
	bw.WriteUInt(0b101, 3)

	result := bw.Bytes()
	expected := []byte{0b10100000}

	if !bytes.Equal(result, expected) {
		t.Fatalf("expected %08b, got %08b", expected, result)
	}
}

func TestZeroSizeWrite(t *testing.T) {
	bw := New()
	bw.WriteUInt(0xFF, 0)

	if len(bw.Bytes()) != 0 {
		t.Fatal("expected no bytes written for size 0")
	}
}

func TestTotalBits(t *testing.T) {
	bw := New()
	bw.WriteUInt(0b101, 3)
	bw.WriteUInt(0b11, 2)

	if bw.TotalBits() != 5 {
		t.Fatalf("expected 5 bits, got %d", bw.TotalBits())
	}
}

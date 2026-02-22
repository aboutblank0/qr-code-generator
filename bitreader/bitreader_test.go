package bitreader

import "testing"

func TestPopSingleByte(t *testing.T) {
	r := New([]byte{0b10101010})

	tests := []struct {
		want bool
	}{
		{false}, // bit 0
		{true},  // bit 1
		{false}, // bit 2
		{true},  // bit 3
		{false}, // bit 4
		{true},  // bit 5
		{false}, // bit 6
		{true},  // bit 7
	}

	for i, tt := range tests {
		got := r.Pop()
		if got != tt.want {
			t.Fatalf("bit %d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestPopAllZeros(t *testing.T) {
	r := New([]byte{0x00})

	for i := range 8 {
		if r.Pop() {
			t.Fatalf("bit %d: expected false, got true", i)
		}
	}
}

func TestPopAllOnes(t *testing.T) {
	r := New([]byte{0xFF})

	for i := range 8 {
		if !r.Pop() {
			t.Fatalf("bit %d: expected true, got false", i)
		}
	}
}

func TestPopMultipleBytes(t *testing.T) {
	r := New([]byte{
		0b00000001,
		0b00000010,
	})

	expected := []bool{
		true, false, false, false, false, false, false, false, // first byte
		false, true, false, false, false, false, false, false, // second byte
	}

	for i, want := range expected {
		got := r.Pop()
		if got != want {
			t.Fatalf("bit %d: got %v, want %v", i, got, want)
		}
	}
}

package gf256

import (
	"testing"
)

func TestTables(t *testing.T) {
	for i := range 255 {
		if int(gfLog[gfExp[i]]) != i {
			t.Fatalf("log/exp mismatch at %d", i)
		}
	}
}

func TestAdd(t *testing.T) {
	cases := []struct{ a, b, want byte }{
		{0b00000000, 0b00000000, 0b0000000},
		{0b00000101, 0b00000101, 0b0000000},
		{0b00000011, 0b00000101, 0b0000110},
	}
	for _, c := range cases {
		if got := Add(c.a, c.b); got != c.want {
			t.Errorf("Add(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestMultiply(t *testing.T) {
	cases := []struct{ a, b, want byte }{
		{0b00000110, 0b00000110, 0b00010100},
		{0b00000010, 0b00000100, 0b00001000},
		{10, 11, 78},
		{87, 19, 224},
		{100, 100, 169},
	}
	for _, c := range cases {
		if got := Multiply(c.a, c.b); got != c.want {
			t.Errorf("Multiply(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestDivide(t *testing.T) {
	cases := []struct {
		a, b, want byte
	}{
		{87, 19, 255},
		{200, 199, 120},
	}
	for _, c := range cases {
		if got := Divide(c.a, c.b); got != c.want {
			t.Errorf("Divide(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestLogExpInverse(t *testing.T) {
	for i := 1; i < 256; i++ {
		a := byte(i)

		// Exp(Log(a)) should recover a (for a ≠ 0)
		exp := Exp(Log(a))
		if exp != a {
			t.Errorf("Exp(Log(%#x)) = %#x, want %#x", a, exp, a)
		}

		// Log(Exp(i)) should give something ≡ i mod 255
		log := Log(Exp(byte(i)))
		expected := byte(i % 255)
		if log != expected {
			t.Errorf("Log(Exp(%d)) = %d, want %d  (for value %#x)", i, log, expected, Exp(byte(i)))
		}
	}
}

func TestMultiplyBy2(t *testing.T) {
	// Just a small table check for x2 multiplication (very common in RS)
	for i := range 256 {
		a := byte(i)
		got := Multiply(a, 2)
		var want byte
		if a&0x80 == 0 {
			want = a << 1
		} else {
			want = (a << 1) ^ 0x1d
		}
		if got != want {
			t.Errorf("2 * %#x = %#x, want %#x", a, got, want)
		}

	}
}

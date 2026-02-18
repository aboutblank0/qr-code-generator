package qr

import (
    "aboutblank/qr-code/gf256"
)

func PolyMultiply(a, b []byte) []byte {
    result := make([]byte, len(a)+len(b)-1)

    for i := range a {
        for j := range b {
            result[i+j] = gf256.Add(result[i+j], gf256.Multiply(a[i], b[j]))
        }
    }

    return result
}

//How to Create a Generator Polynomial (From Thonky.com)
// In each step of creating a generator polynomial, you multiply a polynomial by a polynomial. 
//The very first polynomial that you start with in the first step is always (ɑ0x1 + ɑ0x0).
//For each multiplication step, you multiply the current polynomial by (ɑ^0x^1 + ɑ^j^x0) where j is 1 for the first multiplication, 2 for the second multiplication, 3 for the third, and so on.

func BuildGenerator(r int)[]byte {
    if r < 1 {
	panic("cannot make generator polynomial for r < 1")
    }
    
    curr := []byte{gf256.Exp(0), gf256.Exp(0)}
    for j := 1; j < r; j++ {
	term := []byte{gf256.Exp(0), gf256.Exp(byte(j))}
	curr = PolyMultiply(curr, term)
    }
    return curr
}

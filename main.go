package main

import (
	"aboutblank/qr-code/gf256"
	"aboutblank/qr-code/qr"
	"fmt"
)

func main() {
	//	bytes := qr.WriteAlphanumeric("HELLO WORLD", qr.EC_Q)
	//	fmt.Printf("%08b\n", bytes)

	r := 7
	test := qr.BuildGenerator(r)
	fmt.Printf("r: %d\n", r)
	for _, a := range test {
		fmt.Printf("%d ", gf256.Log(a))
	}
	fmt.Println()

	//res := qr.PolyMultiply([]byte{1, 0, 0, 1}, []byte{1, 0, 0, 1})
	//fmt.Printf("%d\n", res)
	//mult := gf256.Multiply(0b00000001, 0b00000001)
	//fmt.Printf("%08b\n", mult)
}

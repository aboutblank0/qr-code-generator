package main

import (
	"aboutblank/qr-code/qr"
	"fmt"
)

func main() {

	//DELETE ME
	ecBytes := qr.GenerateQRCode("HELLO WORLD", qr.Alphanumeric, qr.EC_M)
	fmt.Printf("%d\n", ecBytes)

	//res := qr.PolyMultiply([]byte{1, 0, 0, 1}, []byte{1, 0, 0, 1})
	//fmt.Printf("%d\n", res)
	//mult := gf256.Multiply(0b00000001, 0b00000001)
	//fmt.Printf("%08b\n", mult)
}

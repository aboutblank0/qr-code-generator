package main

import (
	"aboutblank/qr-code/qr"
	"fmt"
)


func main() {
	bytes := qr.WriteAlphanumeric("HELLO WORLD", qr.EC_L)
	fmt.Printf("%08b\n", bytes)
}


package main

import (
	"aboutblank/qr-code/qr"
	"fmt"
)


func main() {
	bytes := qr.WriteAlphanumeric("Hello")
	fmt.Printf("%08b\n", bytes)
}


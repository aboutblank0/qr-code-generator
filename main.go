package main

import (
	"aboutblank/qr-code/bitwriter"
	"fmt"
)

type EncodingMode uint8

const (
	Numeric EncodingMode = 1 << iota
	Alphanumeric
	Byte
	Kanji
)

const EncodingModeSize uint8 = 4
const TerminatorSize uint8 = 4

func main() {
	bytes := writeQrCode(Alphanumeric, "Hello")
	fmt.Printf("%08b\n", bytes)
}

func writeQrCode(mode EncodingMode, input string) []byte {
	writer := bitwriter.New()

	// Write the encoding mode
	writer.WriteUInt(uint64(mode), EncodingModeSize)

	var version uint8 = 2

	// Write the character count indicator
	//Check which version of QR code we are writing, that defines how many bits (the size) of the character count indicator
	
	charCount := 10 //Of course actually calculate this.

	charCountSize := getCharCountSize(mode, version)
	writer.WriteUInt(uint64(charCount), uint8(charCountSize))

	//Write the actual data

	//Write the terminator
	writer.WriteUInt(0, TerminatorSize)
	return writer.Bytes()
}

func getCharCountSize(m EncodingMode, version uint8) int {
	var group int

	switch {
	case version <= 9:
		group = 0
	case version <= 26:
		group = 1
	default:
		group = 2
	}

	switch group {
	case 0: // Versions 1–9

		switch m {

		case Numeric:
			return 10
		case Alphanumeric:
			return 9
		case Byte:
			return 8
		case Kanji:
			return 8
		}

	case 1: // Versions 10–26
		switch m {
		case Numeric:
			return 12
		case Alphanumeric:
			return 11
		case Byte:
			return 16
		case Kanji:
			return 10
		}

	case 2: // Versions 27–40
		switch m {
		case Numeric:
			return 14
		case Alphanumeric:
			return 13
		case Byte:
			return 16
		case Kanji:
			return 12
		}
	}

	panic("invalid encoding mode or version")
}

package qr

import "aboutblank/qr-code/bitwriter"

type EncodingMode uint8

const (
	Numeric EncodingMode = iota
	Alphanumeric
	Byte
	Kanji
)

const EncodingModeSize uint8 = 4
const TerminatorSize uint8 = 4

func WriteAlphanumeric(input string) []byte {
	writer := bitwriter.New()

	//TODO: Move this to a function to better explain it.
	// Numeric Mode	0001
	//Alphanumeric Mode	0010
	//Byte Mode	0100
	//Kanji Mode	1000
	writer.WriteUInt(uint64(1<<Alphanumeric), EncodingModeSize)

	//TODO: Actually calculate the minimum version required
	var version Version = 2

	// Write the character count indicator
	//Check which version of QR code we are writing, that defines how many bits (the size) of the character count indicator

	charCount := 10 //Of course actually calculate this.

	charCountSize := getCharCountSize(version, Alphanumeric)
	writer.WriteUInt(uint64(charCount), uint8(charCountSize))

	//Write the actual data

	//Write the terminator
	writer.WriteUInt(0, TerminatorSize)
	return writer.Bytes()
}


package qr

import (
	"aboutblank/qr-code/bitwriter"
)

type EncodingMode uint8

const (
	Numeric EncodingMode = iota
	Alphanumeric
	Byte
	Kanji
)

const EncodingModeSize uint8 = 4

func WriteAlphanumeric(input string, ecLevel ErrorCorrectionLevel) []byte {
	writer := bitwriter.New()

	// Write the encoding mode indicator
	writer.WriteUInt(getEncodingModeValue(Alphanumeric), EncodingModeSize)

	// Check if input can be encoded using Alphanumeric AND get the char count.
	charCount, err := getAlphanumericCharCount(input)
	if err != nil {
		panic(err)
	}

	// Determine the QR Code Version
	version, err := determineMinQRVersion(charCount, ecLevel, Alphanumeric)
	if err != nil {
		panic(err)
	}

	// Check which version of QR code we are writing, that defines how many bits (the size) of the character count indicator
	// Write the character count indicator
	charCountSize := getCharCountSize(version, Alphanumeric)
	writer.WriteUInt(uint64(charCount), uint8(charCountSize))

	// Write/Encode the input string
	err = writeAlphanumericString(writer, input)
	if err != nil {
		panic(err)
	}

	ecInfo := getEcInfo(version, ecLevel)
	var requiredBits int = ecInfo.TotalDataCodewords * 8

	//Calculate the amount of terminator bits needed. Maximum of 4 as per the spec
	terminatorSize := requiredBits - writer.TotalBits()
	if terminatorSize > 4 {
		terminatorSize = 4
	}
	writer.WriteUInt(0, uint8(terminatorSize))

	// 2. Pad to next byte boundary
	bitsInLastByte := writer.TotalBits() % 8
	if bitsInLastByte != 0 {

		writer.WriteUInt(0, uint8(8-bitsInLastByte))
	}

	//Add padding bits if still not enough (as per spec)
	remainingBytes := (requiredBits - writer.TotalBits()) / 8
	padBytes := []uint8{0xEC, 0x11}
	for i := 0; i < remainingBytes; i++ {
		writer.WriteUInt(uint64(padBytes[i%2]), 8)
	}

	return writer.Bytes()
}

/*
Numeric			0001
Alphanumeric 	0010
Byte			0100
Kanji			1000
*/
func getEncodingModeValue(mode EncodingMode) uint64 {
	return uint64(1 << mode)
}

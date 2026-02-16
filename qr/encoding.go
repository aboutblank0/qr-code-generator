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
	bitsWritten, err := writeAlphanumericString(writer, input)
	if err != nil {
		panic(err)
	}
	
	ecInfo := getEcInfo(version, ecLevel)
	requiredBits := ecInfo.TotalDataCodewords * 8

	//Calculate the amount of terminator bits needed. Maximum of 4 as per the spec
	terminatorSize := int(requiredBits) - bitsWritten
	if terminatorSize > int(4) {
		terminatorSize = 4
	}
	writer.WriteUInt(0, uint8(terminatorSize))

	//Add padding bits if still not enough

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

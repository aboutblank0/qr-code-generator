package qr

import (
	"aboutblank/qr-code/bitwriter"
	"fmt"
)

type EncodingMode uint8

const (
	Numeric EncodingMode = iota
	Alphanumeric
	Byte
	Kanji
)

const EncodingModeSize uint8 = 4

func GenerateQRCode(input string, encodingMode EncodingMode, ecLevel ErrorCorrectionLevel) []byte {
	writer := bitwriter.New()

	// Write the encoding mode indicator
	writer.WriteUInt(getEncodingModeValue(encodingMode), EncodingModeSize)

	charCount, err := getCharCount(encodingMode, input)
	if err != nil {
		panic(err)
	}

	// Determine the QR Code Version
	version, err := determineMinQRVersion(charCount, ecLevel, encodingMode)
	if err != nil {
		panic(err)
	}

	// Check which version of QR code we are writing, that defines how many bits (the size) of the character count indicator
	// Write the character count indicator
	charCountSize := getCharCountSize(version, encodingMode)
	writer.WriteUInt(uint64(charCount), uint8(charCountSize))

	// Write/Encode the input string
	err = writeAlphanumericString(writer, input)
	if err != nil {
		panic(err)
	}

	ecInfo := getEcInfo(version, ecLevel)
	requiredBits := ecInfo.TotalRequiredBits()

	//Calculate the amount of terminator bits needed. Maximum of 4 as per the spec
	terminatorSize := min(requiredBits-writer.TotalBits(), 4)
	writer.WriteUInt(0, uint8(terminatorSize))

	// Pad to next byte boundary
	bitsInLastByte := writer.TotalBits() % 8
	if bitsInLastByte != 0 {
		writer.WriteUInt(0, uint8(8-bitsInLastByte))
	}

	// Add padding bits if still not enough (as per spec)
	remainingBytes := (requiredBits - writer.TotalBits()) / 8
	padBytes := []uint8{0xEC, 0x11}
	for i := range remainingBytes {
		writer.WriteUInt(uint64(padBytes[i%2]), 8)
	}

	dataCodeWords := writer.Bytes()
	ecCodeWords := GenerateErrorCorrectionCodeWords(dataCodeWords, ecInfo)

	getFinalMessage(dataCodeWords, ecCodeWords, ecInfo)
	return ecCodeWords
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

// Different encoding modes count chars differently
// TODO: Add support for other encoding modes
func getCharCount(mode EncodingMode, data string) (int, error) {
	switch mode{
	case Alphanumeric:
		return getAlphanumericCharCount(data)
	default:
		return 0, fmt.Errorf("invalid encoding mode")
	}
}

// Determines the minimum QR Code version required to "fit" all of the data (charCount)
func determineMinQRVersion(charCount int, ecLevel ErrorCorrectionLevel, mode EncodingMode) (Version, error) {
	for version := Version(1); version <= 40; version++ {
		if capacity[version][ecLevel][mode] >= charCount {
			return version, nil
		}
	}

	return 0, fmt.Errorf("data too long for any QR code version with this ErrorCorrection level")
}

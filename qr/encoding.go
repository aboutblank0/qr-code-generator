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
	requiredBits := ecInfo.TotalDataBits()

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
	fmt.Printf("Data code words: %d\n", dataCodeWords)

	finalMessage := getFinalMessage(dataCodeWords, ecInfo)
	return finalMessage
}

func getFinalMessage(dataCodeWords []byte, ecInfo ecInfo) []byte {
	// ====== Handle Data Code Words =======
	data1 := make([][]byte, 0, ecInfo.Group1.Blocks)
	for i := range cap(data1) {
		block := dataCodeWords[i*ecInfo.Group1.DataCodewords : i*ecInfo.Group1.DataCodewords+ecInfo.Group1.DataCodewords]
		data1 = append(data1, block)
	}

	var data2 [][]byte
	if ecInfo.Group2.Blocks > 0 {
		offset := ecInfo.Group1.Blocks * ecInfo.Group1.DataCodewords
		data2 = make([][]byte, 0, ecInfo.Group2.Blocks)
		for i := range cap(data2) {
			start := offset + i*ecInfo.Group2.DataCodewords
			end := start + ecInfo.Group2.DataCodewords
			block := dataCodeWords[start:end]
			data2 = append(data2, block)
		}
	}

	// ====== Handle Error Correction Code Words =======
	ec1 := make([][]byte, 0, ecInfo.Group1.Blocks)
	for _, data := range data1 {
		ecCodeWords := GenerateErrorCorrectionCodeWords(data, ecInfo)
		fmt.Printf("Error Correction codewords: %d\n", ecCodeWords)
		ec1 = append(ec1, ecCodeWords)
	}

	var ec2 [][]byte
	if data2 != nil {
		ec2 = make([][]byte, 0, ecInfo.Group2.Blocks)
		for _, data := range data2 {
			ecCodeWords := GenerateErrorCorrectionCodeWords(data, ecInfo)
			ec2 = append(ec2, ecCodeWords)
		}
	}

	// ====== Interleave the data Code Words =======
	allDataBlocks := append(data1, data2...)

	maxDataLen := 0
	for _, block := range allDataBlocks {
		if len(block) > maxDataLen {
			maxDataLen = len(block)
		}
	}

	interleavedData := make([]byte, 0, len(allDataBlocks)*maxDataLen)
	for i := 0; i < maxDataLen; i++ {
		for _, block := range allDataBlocks {
			if i < len(block) {
				interleavedData = append(interleavedData, block[i])
			}
		}

	}

	// ====== Interleave Error Correction Codewords ======

	// Combine all EC blocks from both groups
	allECBlocks := append(ec1, ec2...)

	// Find the maximum EC block length
	maxECLen := 0
	for _, block := range allECBlocks {

		if len(block) > maxECLen {
			maxECLen = len(block)
		}
	}

	interleavedEC := make([]byte, 0, len(allECBlocks)*maxECLen)
	for i := 0; i < maxECLen; i++ {

		for _, block := range allECBlocks {

			if i < len(block) {
				interleavedEC = append(interleavedEC, block[i])
			}
		}
	}

	finalCodewords := append(interleavedData, interleavedEC...)
	return finalCodewords
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

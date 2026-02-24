package qr

import (
	"aboutblank/qr-code/bitwriter"
	"fmt"
)

type EncodingMode uint8

const (
	Encode_Numeric EncodingMode = iota
	Encode_Alphanumeric
	Encode_Byte
	Encode_Kanji
)

/*
Numeric			0001
Alphanumeric 	0010
Byte			0100
Kanji			1000
*/
func getEncodingModeValue(mode EncodingMode) uint64 {
	return uint64(1 << mode)
}

func getEncodingModeString(mode EncodingMode) string {
	switch mode {
	case Encode_Numeric:
		return "Numeric"
	case Encode_Alphanumeric:
		return "Alphanumeric"
	case Encode_Kanji:
		return "Kanji"
	case Encode_Byte:
		return "Byte"
	}
	return "INVALID"
}

func getErrorCorrectionString(ecLevel ErrorCorrectionLevel) string {
	switch ecLevel {
	case EC_Low:
		return "Low"
	case EC_Medium:
		return "Medium"
	case EC_Quartile:
		return "Quartile"
	case EC_High:
		return "High"
	}
	return "INVALID"
}

func GenerateQRCode(input string, ecLevel ErrorCorrectionLevel) *QRCode {
	writer := bitwriter.New()

	encodingMode := determineBestEncodingMode(input)
	fmt.Printf("Encoding Mode: %s\n", getEncodingModeString(encodingMode))
	fmt.Printf("Error Correction Level: %s\n", getErrorCorrectionString(ecLevel))

	// Write the encoding mode indicator (always 4 bits)
	writer.WriteUInt(getEncodingModeValue(encodingMode), 4)

	charCount, err := getCharCount(encodingMode, input)
	if err != nil {
		panic(err)
	}

	// Determine the QR Code Version
	version, err := determineMinQRVersion(charCount, ecLevel, encodingMode)
	if err != nil {
		panic(err)
	}
	fmt.Printf("QRCode Version: %d\n", version)


	// Check which version of QR code we are writing, that defines how many bits (the size) of the character count indicator
	// Write the character count indicator
	charCountSize := getCharCountSize(version, encodingMode)
	writer.WriteUInt(uint64(charCount), uint8(charCountSize))
	fmt.Printf("Writing char count: %d\n", charCount)

	// Write/Encode the input string
	err = writeString(writer, encodingMode, input)
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

	qrCode := New(version, ecLevel)
	qrCode.ApplyFinalMessage(finalMessage)
	return qrCode
}

func getFinalMessage(dataCodeWords []byte, ecInfo ErrorCorrectionInfo) []byte {
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
		ecCodeWords := generateErrorCorrectionCodeWords(data, ecInfo)
		fmt.Printf("Error Correction code words: %d\n", ecCodeWords)
		ec1 = append(ec1, ecCodeWords)
	}

	var ec2 [][]byte
	if data2 != nil {
		ec2 = make([][]byte, 0, ecInfo.Group2.Blocks)
		for _, data := range data2 {
			ecCodeWords := generateErrorCorrectionCodeWords(data, ecInfo)
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

// Different encoding modes count chars differently
func getCharCount(mode EncodingMode, data string) (int, error) {
	switch mode {
	case Encode_Numeric:
		return getNumericCharCount(data)
	case Encode_Alphanumeric:
		return getAlphanumericCharCount(data)
	case Encode_Byte:
		return getByteCharCount(data)
	case Encode_Kanji:
		return getKanjiCharCount(data)
	default:
		return 0, fmt.Errorf("invalid encoding mode")
	}
}

func writeString(writer *bitwriter.BitWriter, mode EncodingMode, data string) error {
	var err error

	switch mode {
	case Encode_Numeric:
		err = writeNumericString(writer, data)
	case Encode_Alphanumeric:
		err = writeAlphanumericString(writer, data)
	case Encode_Byte:
		err = writeByteString(writer, data)
	case Encode_Kanji:
		err = writeKanjiString(writer, data)
	default:
		return fmt.Errorf("invalid encoding mode")
	}

	return err
}

// Determines the minimum QR Code version required to "fit" all of the data (charCount)
func determineMinQRVersion(charCount int, ecLevel ErrorCorrectionLevel, mode EncodingMode) (Version, error) {
	for version := Version(1); version <= 40; version++ {
		if getMaxCharCapacity(version, ecLevel, mode) >= charCount {
			return version, nil
		}
	}

	return 0, fmt.Errorf("data too long for any QR code version with this ErrorCorrection level")
}

func determineBestEncodingMode(data string) EncodingMode {
	if canEncodeNumeric(data) {
		return Encode_Numeric
	}

	if canEcodeAlphanumeric(data) {
		return Encode_Alphanumeric
	}
	
	if canEncodeKanji(data) {
		return Encode_Kanji
	}

	return Encode_Byte
}

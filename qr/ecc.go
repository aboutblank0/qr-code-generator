package qr

import "aboutblank/qr-code/gf256"

type ErrorCorrectionLevel int

const (
	EC_Low ErrorCorrectionLevel = iota
	EC_Medium
	EC_Quartile
	EC_High
)

type blockGroup struct {
	Blocks        int // number of blocks
	DataCodewords int // data codewords per block
}

type ErrorCorrectionInfo struct {
	TotalDataCodewords  int
	ECCodewordsPerBlock int
	Group1              blockGroup
	Group2              blockGroup // Blocks == 0 if unused
}

func generateErrorCorrectionCodeWords(dataCodeWords []byte, ecInfo ErrorCorrectionInfo) []byte {
	genPoly := buildGenerator(ecInfo.ECCodewordsPerBlock)

	// Multiply message by x^n 
	messagePoly := shiftPoly(dataCodeWords, ecInfo.ECCodewordsPerBlock)

	// Division loop
	for i := 0; i <= len(messagePoly)-len(genPoly); i++ {
		lead := messagePoly[i]
		if lead == 0 {
			continue
		}
		for j := range genPoly {
			messagePoly[i+j] ^= gf256.Multiply(genPoly[j], lead)
		}
	}

	// Last n bytes = EC codewords
	ecBytes := messagePoly[len(messagePoly)-ecInfo.ECCodewordsPerBlock:]
	return ecBytes
}

func getEcInfo(version Version, ecLevel ErrorCorrectionLevel) ErrorCorrectionInfo {
	return ecTable[version][ecLevel]
}

func (e ErrorCorrectionInfo) TotalBlocks() int {
	return e.Group1.Blocks + e.Group2.Blocks
}

func (e ErrorCorrectionInfo) TotalECCodewords() int {
	return e.TotalBlocks() * e.ECCodewordsPerBlock
}

func (e ErrorCorrectionInfo) TotalCodewords() int {
	return e.TotalDataCodewords + e.TotalECCodewords()
}

func (e ErrorCorrectionInfo) TotalDataBits() int {
	return e.TotalDataCodewords * 8
}

func (e ErrorCorrectionInfo) TotalRequiredBits() int {
	return e.TotalCodewords() * 8
}

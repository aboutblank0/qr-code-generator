package qr

type Version uint8

func getFinalMessage(dataCodeWords, _ []byte, ecInfo ecInfo) []byte {
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

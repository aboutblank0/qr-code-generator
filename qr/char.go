package qr

// getMaxCharCapacity() returns the maximum amount of characters
// for the given QR version, error correction level and encoding mode.
func getMaxCharCapacity(version Version, ecLevel ErrorCorrectionLevel, mode EncodingMode) int {
	return maxCharCapacity[version][ecLevel][mode]
}

// getCharCountSize returns the bit length of the character count indicator
// for the given QR version and encoding mode.
//
// The size depends on the version group:
// 1–9, 10–26, or 27–40.
func getCharCountSize(version Version, mode EncodingMode) int {
	var group int

	switch {
	case version <= 9:
		group = 0
	case version <= 26:
		group = 1
	default:
		group = 2
	}

	return charCountSize[group][mode]
}

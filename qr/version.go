package qr

import "fmt"

type Version uint8

func determineMinQRVersion(charCount int, ecLevel ErrorCorrectionLevel, mode EncodingMode) (Version, error) {
	for version := Version(1); version <= 40; version++ {
		if capacity[version][ecLevel][mode] >= charCount {
			return version, nil
		}
	}

	return 0, fmt.Errorf("data too long for any QR code version with this ErrorCorrection level")
}

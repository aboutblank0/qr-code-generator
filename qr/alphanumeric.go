package qr

import (
	"aboutblank/qr-code/bitwriter"
	"fmt"
)

var alphaNumMap = [128]uint8{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4,
	'5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16,
	'H': 17, 'I': 18, 'J': 19, 'K': 20, 'L': 21, 'M': 22, 'N': 23,
	'O': 24, 'P': 25, 'Q': 26, 'R': 27, 'S': 28, 'T': 29, 'U': 30,
	'V': 31, 'W': 32, 'X': 33, 'Y': 34, 'Z': 35,
	' ': 36, '$': 37, '%': 38, '*': 39, '+': 40,
	'-': 41, '.': 42, '/': 43, ':': 44,
}

func alphanumericValue(r rune) (uint, bool) {
	if r > 127 || alphaNumMap[r] == 0 && r != '0' {
		return 0, false
	}
	return uint(alphaNumMap[r]), true
}

func getAlphanumericCharCount(s string) (int, error) {
	count := 0
	for i, r := range s {
		if _, ok := alphanumericValue(r); !ok {
			return 0, fmt.Errorf("invalid character at position %d: %q", i, r)
		}
		count++
	}
	return count, nil
}


func writeAlphanumericString(writer *bitwriter.BitWriter, s string) (int, error) {
	runes := []rune(s)
	bitsWritten := 0
	for i := 0; i < len(runes); i += 2 {
		if i+1 < len(runes) {
			ch1, ok := alphanumericValue(runes[i])
			ch2, ok2 := alphanumericValue(runes[i+1])

			if !ok || !ok2 {
				return 0, fmt.Errorf("invalid character for alphanumeric encoding %q", runes[i])
			}

			encoded := uint64(ch1) * 45 + uint64(ch2)
			writer.WriteUInt(encoded, 11)
			bitsWritten += 11
		} else {
			ch1, ok := alphanumericValue(runes[i])
			if !ok {
				return 0, fmt.Errorf("invalid character for alphanumeric encoding %q", runes[i])
			}

			writer.WriteUInt(uint64(ch1), 6)
			bitsWritten += 6
		}
	}
	return bitsWritten, nil
}

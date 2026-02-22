package qr

import (
	"aboutblank/qr-code/bitwriter"
	"fmt"
)

var numMap = map[rune]int {
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4,
	'5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
}

func canEncodeNumeric(s string) bool {
	for _, r := range s {
		if _, ok := numMap[r]; !ok {
			return false
		}
	}
	return true
}

func getNumericCharCount(s string) (int, error) {
	count := 0
	for i, r := range s {
		if _, ok := numMap[r]; !ok {
			return 0, fmt.Errorf("invalid character at position %d: %q", i, r)
		}
		count++
	}
	return count, nil
}

func writeNumericString(writer *bitwriter.BitWriter, s string) error {
	runes := []rune(s)
	i := 0

	// process groups of 3
	for i+2 < len(runes) {
		val := 100*numMap[runes[i]] + 10*numMap[runes[i+1]] + numMap[runes[i+2]]
		writer.WriteUInt(uint64(val), 10)
		i += 3
	}

	// handle leftovers
	if i < len(runes) {
		if i+1 < len(runes) {
			// two digits left -> 7 bits
			fmt.Println("handling writing 2 digits")
			val := 10*numMap[runes[i]] + numMap[runes[i+1]]
			writer.WriteUInt(uint64(val), 7)
			return nil
		}

		// one digit left -> 4 bits
		val := numMap[runes[i]]
		writer.WriteUInt(uint64(val), 4)
		return nil
	}

	return nil
}


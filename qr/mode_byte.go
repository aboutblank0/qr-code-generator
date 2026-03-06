package qr

import (
	"aboutblank/qr-code/bitwriter"
)

func getByteCharCount(s string) (int, error) {
	return len([]byte(s)), nil
}

func writeByteString(writer *bitwriter.BitWriter, s string) error {
	data := []byte(s)
	for _, b := range data {
		writer.WriteUInt(uint64(b), 8)
	}
	return nil
}

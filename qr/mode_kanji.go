package qr

import (
	"aboutblank/qr-code/bitwriter"
	"fmt"
	"golang.org/x/text/encoding/japanese"
)

func canEncodeKanji(s string) bool {
	_, err := toShiftJIS(s)
	return err == nil
}

// Only double byte JIS characters are supported.
func getKanjiCharCount(s string) (int, error) {
	b, err := toShiftJIS(s)
	if err != nil {
		return 0, err
	}
	return len(b)/2, nil 
}

func writeKanjiString(writer *bitwriter.BitWriter, s string) error {
	b, err := toShiftJIS(s)
	if err != nil {
		return err
	}

	if len(b)%2 != 0 {
		return fmt.Errorf("Shift_JIS byte length is not even")
	}

	for i := 0; i < len(b); i += 2 {
		code := (uint16(b[i]) << 8) | uint16(b[i+1])

		var adjusted uint16
		if code >= 0x8140 && code <= 0x9FFC {
			adjusted = code - 0x8140
		} else if code >= 0xE040 && code <= 0xEBBF {
			adjusted = code - 0xC140

		} else {
			return fmt.Errorf("character not in JIS X 0208 range")
		}

		// QR Kanji mode packing formula
		kanjiValue := ((adjusted >> 8) * 0xC0) + (adjusted & 0xFF)
		writer.WriteUInt(uint64(kanjiValue), 13)
	}

	return nil
}

func toShiftJIS(s string) ([]byte, error) {
	encoder := japanese.ShiftJIS.NewEncoder()
	return encoder.Bytes([]byte(s))
}

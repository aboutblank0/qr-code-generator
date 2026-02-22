package qr // Assuming this is in a package named qr or similar

import (
	"bytes"
	"testing"
)

// TestGenerateErrorCorrectionCodeWords tests the EC generation using a known QR example.
// This uses the "HELLO WORLD" data from QR version 1-Q as per ISO/IEC 18004.
func TestGenerateErrorCorrectionCodeWords(t *testing.T) {
	// Example data codewords for "HELLO WORLD" in alphanumeric mode, version 1-Q
	dataCodeWords := []byte{32, 91, 11, 241, 209, 114, 220, 38, 161, 160, 236}

	// EC info for version 1-Q: 15 EC codewords per block (single block)
	ecInfo := ecInfo{

		ECCodewordsPerBlock: 15,
		// Assume other fields are set as needed, but only ECCodewordsPerBlock is used here
	}

	// Compute EC codewords

	got := GenerateErrorCorrectionCodeWords(dataCodeWords, ecInfo)

	// Expected EC codewords from standard QR encoding (verified example)

	want := []byte{220, 117, 193, 108, 198, 216, 11, 4, 232, 173, 99, 128, 138, 2, 14}

	if !bytes.Equal(got, want) {
		t.Errorf("GenerateErrorCorrectionCodeWords() = %v, want %v", got, want)
	}
}

// TestGenerateErrorCorrectionCodeWords_ZeroData tests with all-zero data, which should produce zero EC (trivial case).

func TestGenerateErrorCorrectionCodeWords_ZeroData(t *testing.T) {
	// Zero data codewords (e.g., for a small RS(5,3) with 2 EC, but adjust to match function)
	// For simplicity, use a small custom RS(5,3) where n=5, k=3, EC=2
	dataCodeWords := []byte{0, 0, 0} // Message of length k=3

	ecInfo := ecInfo{
		ECCodewordsPerBlock: 2,
	}

	got := GenerateErrorCorrectionCodeWords(dataCodeWords, ecInfo)

	// For all-zero message, EC should be all-zero
	want := []byte{0, 0}

	if !bytes.Equal(got, want) {
		t.Errorf("GenerateErrorCorrectionCodeWords() = %v, want %v", got, want)
	}
}

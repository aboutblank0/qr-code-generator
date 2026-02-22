package qr

// represents the "shape" of the finder patterns found
// on the top-left, top-right and bottom-left corners
// of a QR Code.
var finderPattern = [7][7]bool{
	{true, true,  true,	 true,  true,  true,  true},
	{true, false, false, false, false, false, true},
	{true, false, true,  true,  true,  false, true},
	{true, false, true,  true,  true,  false, true},
	{true, false, true,  true,  true,  false, true},
	{true, false, false, false, false, false, true},
	{true, true,  true,  true,  true,  true,  true},
}

// represents the "shape" of the alignment patterns found
// all spread out through the QR Code for version > 1
var alignmentPattern = [5][5]bool{
	{true, true,  true,  true, true},
	{true, false, false, false, true},
	{true, false, true,  false, true},
	{true, false, false, false, true},
	{true, true,  true,  true, true},
}

// Top-left coordinates of alignment patterns for each version.
// Empty slice means no alignment patterns (version 1).
var alignmentPatternPositions = [41][]int{
	0: {}, 1: {},
	2:  {4, 16},
	3:  {4, 20},
	4:  {4, 24},
	5:  {4, 28},
	6:  {4, 32},
	7:  {4, 20, 36},
	8:  {4, 22, 40},
	9:  {4, 24, 44},
	10: {4, 26, 48},
	11: {4, 28, 52},
	12: {4, 30, 56},
	13: {4, 32, 60},
	14: {4, 24, 44, 64},
	15: {4, 24, 46, 68},
	16: {4, 24, 48, 72},
	17: {4, 28, 52, 76},
	18: {4, 28, 54, 80},
	19: {4, 28, 56, 84},
	20: {4, 32, 60, 88},
	21: {4, 26, 48, 70, 92},
	22: {4, 24, 48, 72, 96},
	23: {4, 28, 52, 76, 100},
	24: {4, 26, 52, 78, 104},
	25: {4, 30, 56, 82, 108},
	26: {4, 28, 56, 84, 112},
	27: {4, 32, 60, 88, 116},
	28: {4, 24, 48, 72, 96, 120},
	29: {4, 28, 52, 76, 100, 124},
	30: {4, 24, 50, 76, 102, 128},
	31: {4, 28, 54, 80, 106, 132},
	32: {4, 32, 58, 84, 108, 136},
	33: {4, 28, 56, 84, 112, 140},
	34: {4, 32, 60, 88, 116, 144},
	35: {4, 28, 52, 76, 100, 124, 148},
	36: {4, 22, 48, 74, 100, 126, 152},
	37: {4, 26, 52, 78, 104, 130, 156},
	38: {4, 30, 56, 82, 108, 134, 160},
	39: {4, 24, 52, 80, 108, 136, 164},
	40: {4, 28, 56, 84, 112, 138, 168},
}

package qr

import (
	"image"
)

type Version uint8

type QRCode struct {
	Version      Version
	ModuleMatrix [][]Module
}

type Module struct {
	Value    bool
	Reserved bool
}

var finderPattern = [7][7]bool{
	{true, true, true, true, true, true, true},
	{true, false, false, false, false, false, true},
	{true, false, true, true, true, false, true},
	{true, false, true, true, true, false, true},
	{true, false, true, true, true, false, true},
	{true, false, false, false, false, false, true},
	{true, true, true, true, true, true, true},
}

var alignmentPattern = [5][5]bool{
	{true, true, true, true, true},
	{true, false, false, false, true},
	{true, false, true, false, true},
	{true, false, false, false, true},
	{true, true, true, true, true},
}

// Top-left coordinates of alignment patterns for each version.
// Empty slice means no alignment patterns (version 1).
var alignmentsTopLeft = [41][]int{
	0:  {},
	1:  {},
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

func New(version Version) *QRCode {
	qr := &QRCode{}
	qr.Version = version

	size := qr.GetSize()
	qr.ModuleMatrix = make([][]Module, size)
	for i := range qr.ModuleMatrix {
		qr.ModuleMatrix[i] = make([]Module, size)
	}
	return qr
}

func (qr QRCode) GetSize() int {
	return int(21 + 4*(qr.Version-1))
}

func (qr *QRCode) getModule(x, y int) Module {
	return qr.ModuleMatrix[x][y]
}

func (qr *QRCode) setModule(x, y int, value, reserved bool) {
	qr.ModuleMatrix[x][y].Value = value
	qr.ModuleMatrix[x][y].Reserved = reserved
}

func (qr *QRCode) AddFinderPatternsAndSeparators() {
	size := qr.GetSize()

	qr.addFinder(0, 0) // Top left
	qr.addSeparatorLine(7, 0, true)
	qr.addSeparatorLine(0, 7, false)

	qr.addFinder(0, size-7) // Bottom left
	qr.addSeparatorLine(7, size-8, true)
	qr.addSeparatorLine(0, size-8, false)

	qr.addFinder(size-7, 0) // Top right
	qr.addSeparatorLine(size-8, 0, true)
	qr.addSeparatorLine(size-8, 7, false)
}

func (qr *QRCode) AddAlignmentPatterns() {
	if qr.Version < 2 {
		return
	}

	positions := alignmentsTopLeft[qr.Version]
	for _, pos := range positions {
		qr.addAlignmentPattern(pos, pos)
	}
}

func (qr *QRCode) addFinder(x0, y0 int) {
	for x := range 7 {
		for y := range 7 {
			qr.setModule(x0+x, y0+y, finderPattern[x][y], true)
		}
	}
}

// Length is always 8
// If !vertical then obviously horizontaal
func (qr *QRCode) addSeparatorLine(x0, y0 int, vertical bool) {
	length := 8

	for i := range length {
		if vertical {
			qr.setModule(x0, y0+i, false, true)
		} else {
			qr.setModule(x0+i, y0, false, true)
		}
	}
}

// Returns whether the pattern was added or not.
// False if it was "obstructed"
func (qr *QRCode) addAlignmentPattern(x0, y0 int) bool {
	// Check whether any of the positions are already reserved
	for x := range 5 {
		for y := range 5 {
			if qr.getModule(x0+x, y0+y).Reserved {
				return false
			}
		}
	}

	// Actually "write" the alignment pattern
	for x := range 5 {
		for y := range 5 {
			qr.setModule(x0+x, y0+y, alignmentPattern[x][y], true)
		}
	}

	return true
}

func (qr *QRCode) GenerateImage(scale int) *image.RGBA {
	size := qr.GetSize()
	w, h := size*scale, size*scale
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	pix := img.Pix
	stride := img.Stride

	for x := range size {
		for y := range size {
			c := byte(255) // white
			a := byte(0)
			if qr.ModuleMatrix[x][y].Value {
				c = 0
			}

			if qr.ModuleMatrix[x][y].Reserved {
				a = 255
			}

			// Fill the scale×scale block
			for dy := range scale {
				rowStart := (y*scale+dy)*stride + x*scale*4
				for dx := range scale {
					offset := rowStart + dx*4
					pix[offset+0] = c // R
					pix[offset+1] = c // G
					pix[offset+2] = c // B
					pix[offset+3] = a // A
				}
			}
		}
	}

	return img
}

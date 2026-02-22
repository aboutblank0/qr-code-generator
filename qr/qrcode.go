package qr

import (
	"aboutblank/qr-code/bitreader"
	"image"
)

type Version uint8

type QRCode struct {
	Version      Version
	ModuleMatrix [][]Module
}

type Module struct {
	Value    ModuleValue
	Reserved bool
}

type ModuleValue uint8

const (
	ValueNone ModuleValue = iota
	ValueBlack
	ValueWhite
)

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

func (qr *QRCode) setModule(x, y int, value ModuleValue, reserved bool) {
	qr.ModuleMatrix[x][y].Value = value
	qr.ModuleMatrix[x][y].Reserved = reserved
}

func (qr *QRCode) Test(data []byte) {
	qr.AddFinderPatternsAndSeparators()
	qr.AddAlignmentPatterns()
	qr.AddTimingPatterns()
	qr.AddDarkModule()

	qr.ReserveFormatModules()
	qr.ReserveVersionModules()
	qr.WriteData(data)
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
	for _, pos1 := range positions {
		for _, pos2 := range positions {
			qr.addAlignmentPattern(pos1, pos2)
		}
	}
}

func (qr *QRCode) AddTimingPatterns() {
	// Vertical timing pattern
	x0, y0 := 6, 6

	val := true
	for y := 0; y < len(qr.ModuleMatrix[x0])-y0; y++ {
		if !qr.getModule(x0, y0+y).Reserved {
			if val {
				qr.setModule(x0, y0+y, ValueBlack, true)
			} else {
				qr.setModule(x0, y0+y, ValueWhite, true)
			}
		}
		val = !val
	}

	// Horizontal timing pattern
	val = true
	for x := 0; x < len(qr.ModuleMatrix)-x0; x++ {
		if !qr.getModule(x0+x, y0).Reserved {
			if val {
				qr.setModule(x0+x, y0, ValueBlack, true)
			} else {
				qr.setModule(x0+x, y0, ValueWhite, true)
			}
		}
		val = !val
	}
}

func (qr *QRCode) AddDarkModule() {
	//(8, [(4 * V) + 9])
	x := 8
	y := int((4 * qr.Version) + 9)

	qr.setModule(x, y, ValueBlack, true)
}

func (qr *QRCode) ReserveFormatModules() {
	// Top Left (vertical)
	x := 8
	y := 0
	for i := range 8 {
		if !qr.getModule(x, y+i).Reserved {
			qr.setModule(x, y+i, ValueNone, true)
		}
	}

	// Top Left (horizontal)
	x = 0
	y = 8
	for i := range 9 {
		if !qr.getModule(x+i, y).Reserved {
			qr.setModule(x+i, y, ValueNone, true)
		}
	}

	// Top right
	x = qr.GetSize() - 8
	y = 8
	for i := range 8 {
		qr.setModule(x+i, y, ValueNone, true)
	}

	// Bottom Left
	x = 8
	y = int((4 * qr.Version) + 9)
	for i := range 8 {
		if !qr.getModule(x, y+i).Reserved {
			qr.setModule(x, y+i, ValueNone, true)
		}
	}
}

// Only applies to QRCode version 7 and above
func (qr *QRCode) ReserveVersionModules() {
	if qr.Version < 7 {
		return
	}

	// Top right
	startX := qr.GetSize() - 11
	startY := 0
	for x := range 3 {
		for y := range 6 {
			qr.setModule(startX+x, startY+y, ValueNone, true)
		}
	}

	// Bottom left
	startX = 0
	startY = qr.GetSize() - 11
	for x := range 6 {
		for y := range 3 {
			qr.setModule(startX+x, startY+y, ValueNone, true)
		}
	}
}

func (qr *QRCode) WriteData(data []byte) {
	size := qr.GetSize()
	skipX := 6 //x coordinate to "skip"

	reader := bitreader.New(data)

	x, y := size-1, size-1


	// Draw first module
	if !qr.getModule(x, y).Reserved {
		if reader.Pop() {
			qr.setModule(x, y, ValueBlack, false)
		} else {
			qr.setModule(x, y, ValueWhite, false)
		}
	}

	yDir := -1
	moveY := false
	wrapped := true

	for reader.HasData() {
		/*
			// Handle wrapping
			if y == 0 || y == size-1 {
				movedY = true
				yDir *= -1
				x -= 2
			}
		*/

		// Handle zig zag
		if moveY {
			if (y == 0 || y == size-1) && !wrapped { // Handle Wrapping
				wrapped = true
				x--

				if x == skipX { // Special case to skip over timer pattern
					x--
				}

				yDir *= -1
			} else {
				x++
				y += yDir
				wrapped = false
			}
			moveY = false
		} else {
			x--
			moveY = true
		}

		if !qr.getModule(x, y).Reserved {
			if reader.Pop() {
				qr.setModule(x, y, ValueBlack, false)
			} else {
				qr.setModule(x, y, ValueWhite, false)
			}
		}
	}
}

func (qr *QRCode) addFinder(x0, y0 int) {
	for x := range 7 {
		for y := range 7 {
			if finderPattern[x][y] {
				qr.setModule(x0+x, y0+y, ValueBlack, true)
			} else {
				qr.setModule(x0+x, y0+y, ValueWhite, true)
			}
		}
	}
}

// Length is always 8
// If !vertical then obviously horizontaal
func (qr *QRCode) addSeparatorLine(x0, y0 int, vertical bool) {
	length := 8

	for i := range length {
		if vertical {
			qr.setModule(x0, y0+i, ValueWhite, true)
		} else {
			qr.setModule(x0+i, y0, ValueWhite, true)
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
			if alignmentPattern[x][y] {
				qr.setModule(x0+x, y0+y, ValueBlack, true)
			} else {
				qr.setModule(x0+x, y0+y, ValueWhite, true)
			}
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

			switch qr.ModuleMatrix[x][y].Value {
			case ValueBlack:
				c = 0
			case ValueNone:
				c = 123
			}

			// TODO: DELETE. Debugging
			if qr.ModuleMatrix[x][y].Value == ValueNone && qr.ModuleMatrix[x][y].Reserved {
				c = 50
			}

			// Fill the scale×scale block
			for dy := range scale {
				rowStart := (y*scale+dy)*stride + x*scale*4
				for dx := range scale {
					offset := rowStart + dx*4
					pix[offset+0] = c   // R
					pix[offset+1] = c   // G
					pix[offset+2] = c   // B
					pix[offset+3] = 255 // A
				}
			}
		}
	}

	return img
}

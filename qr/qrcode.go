package qr

import (
	"aboutblank/qr-code/bitreader"
	"image"
)

type Version uint8

type QRCode struct {
	Version      Version
	EcLevel      ErrorCorrectionLevel
	EncodingMode EncodingMode

	moduleMatrix [][]Module
	mask         int

	size             int
	formatPositions  [30][2]int
	versionPositions [36][2]int
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


func New(version Version, ecLevel ErrorCorrectionLevel) *QRCode {
	qr := &QRCode{}
	qr.Version = version
	qr.EcLevel = ecLevel
	qr.mask = -1

	size := int(21 + 4*(version-1))
	qr.size = size

	qr.moduleMatrix = make([][]Module, qr.size)
	for i := range qr.moduleMatrix {
		qr.moduleMatrix[i] = make([]Module, qr.size)
	}

	// NOTE: It's important that they are WRITTEN to in this specific order. (format and version)

	qr.formatPositions = [30][2]int{
		// Top-left area
		{0, 8}, {1, 8}, {2, 8}, {3, 8}, {4, 8}, {5, 8}, {7, 8}, {8, 8},
		{8, 7}, {8, 5}, {8, 4}, {8, 3}, {8, 2}, {8, 1}, {8, 0},

		// Bottom left 
		{8, size - 1}, {8, size - 2}, {8, size - 3}, {8, size - 4},{8, size - 5},{8, size - 6}, {8, size - 7},

		// Top Right
		{size - 8, 8}, {size - 7, 8}, {size - 6, 8}, {size - 5, 8}, {size - 4, 8}, {size - 3, 8}, {size - 2, 8}, {size - 1, 8}, 
	}

	qr.versionPositions = [36][2]int{
		// Bottom-left block (6x3)
		{0, size - 11}, {0, size - 10}, {0, size - 9}, {1, size - 11}, {1, size - 10}, {1, size - 9},
		{2, size - 11}, {2, size - 10}, {2, size - 9}, {3, size - 11}, {3, size - 10}, {3, size - 9},
		{4, size - 11}, {4, size - 10}, {4, size - 9}, {5, size - 11}, {5, size - 10}, {5, size - 9},


		// Top-right block (3x6)
		{size - 11, 0}, {size - 10, 0}, {size - 9, 0},
		{size - 11, 1}, {size - 10, 1}, {size - 9, 1},
		{size - 11, 2}, {size - 10, 2}, {size - 9, 2},
		{size - 11, 3}, {size - 10, 3}, {size - 9, 3},
		{size - 11, 4}, {size - 10, 4}, {size - 9, 4},
		{size - 11, 5}, {size - 10, 5}, {size - 9, 5},
	}

	return qr
}

func (qr *QRCode) Clone() *QRCode {
	size := len(qr.moduleMatrix)

	newMatrix := make([][]Module, size)
	for i := range newMatrix {
		newMatrix[i] = make([]Module, size)
		copy(newMatrix[i], qr.moduleMatrix[i])
	}

	return &QRCode{
		Version:      qr.Version,
		moduleMatrix: newMatrix,
	}
}

func (qr *QRCode) getModule(x, y int) Module {
	return qr.moduleMatrix[x][y]
}

func (qr *QRCode) setModule(x, y int, value ModuleValue, reserved bool) {
	qr.moduleMatrix[x][y].Value = value
	qr.moduleMatrix[x][y].Reserved = reserved
}

func (qr *QRCode) ApplyFinalMessage(data []byte) {
	qr.AddFinderPatternsAndSeparators()
	qr.AddAlignmentPatterns()
	qr.AddTimingPatterns()
	qr.AddDarkModule()
	qr.ReserveModules()

	qr.WriteData(data)
	qr.ApplyBestMask()

	qr.WriteFormatInfo()
	qr.WriteVersionInfo()
}

func (qr *QRCode) AddFinderPatternsAndSeparators() {
	size := qr.size

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

	positions := alignmentPatternPositions[qr.Version]
	for _, pos1 := range positions {
		for _, pos2 := range positions {
			qr.addAlignmentPattern(pos1, pos2)
		}
	}
}

func (qr *QRCode) AddTimingPatterns() {
	x0, y0 := 6, 6 // start position for the timing pattern

	// Vertical timing pattern
	val := true
	for y := 0; y < len(qr.moduleMatrix[x0])-y0; y++ {
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
	for x := 0; x < len(qr.moduleMatrix)-x0; x++ {
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

func (qr *QRCode) ReserveModules() {
	for _, p := range qr.formatPositions {
		qr.setModule(p[0], p[1], ValueNone, true)
	}

	if qr.Version >= 7 {
		for _, p := range qr.versionPositions {
			qr.setModule(p[0], p[1], ValueNone, true)
		}
	}
}

func (qr *QRCode) WriteData(data []byte) {
	reader := bitreader.New(data)
	positions := qr.dataPositions()

	for _, pos := range positions {

		val := ValueWhite
		if reader.HasData() && reader.Pop() {
			val = ValueBlack
		}

		qr.setModule(pos[0], pos[1], val, false)
	}
}

func (qr *QRCode) WriteFormatInfo() {
	info := formatInfo[qr.EcLevel][qr.mask]

	for i, pos := range qr.formatPositions {
		bitIndex := 14 - (i % 15)
		bit := (info >> bitIndex) & 1

		val := ValueWhite
		if bit == 1 {
			val = ValueBlack
		}

		qr.setModule(pos[0], pos[1], val, false)
	}
}

func (qr *QRCode) WriteVersionInfo() {
	if qr.Version < 7 {
		return
	}

	info := versionInfo[int(qr.Version)]
	for i, pos := range qr.versionPositions {
		bit := (info >> (17 - i)) & 1
		val := ValueWhite
		if bit == 1 {
			val = ValueBlack
		}
		qr.setModule(pos[0], pos[1], val, false)
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

func (qr *QRCode) dataPositions() [][2]int {
	var positions [][2]int

	size := qr.size
	x := size - 1
	y := size - 1
	dir := -1

	for x > 0 {
	    if x == 6 { // skip timing column
			x--
		}

		for {
			for i := range 2 {
				xx := x - i
				if !qr.getModule(xx, y).Reserved {
					positions = append(positions, [2]int{xx, y})
				}
			}

			y += dir
			if y < 0 || y >= size {
				y -= dir
				dir *= -1
				break
			}
		}
		x -= 2
	}

	return positions
}

func (qr *QRCode) GenerateImage(scale int) *image.RGBA {
	size := qr.size
	padding := 4 // quiet zone (light modules)
	gridSize := size + (padding * 2)

	w, h := gridSize*scale, gridSize*scale
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	pix := img.Pix
	stride := img.Stride

	// make everything white
	for i := range pix {
		pix[i] = 255
	}

	for x := range size {
		for y := range size {
			c := byte(255)

			switch qr.moduleMatrix[x][y].Value {
			case ValueBlack:
				c = 0
			case ValueNone:
				c = 123
			}

			// TODO: REMOVE ME DEBUGGING
			if qr.moduleMatrix[x][y].Value == ValueNone && qr.moduleMatrix[x][y].Reserved {
				c = 50
			}

			drawX := (x + padding) * scale
			drawY := (y + padding) * scale
			for dy := range scale {
				rowStart := (drawY+dy)*stride + drawX*4
				for dx := range scale {
					offset := rowStart + dx*4
					pix[offset+0] = c
					pix[offset+1] = c
					pix[offset+2] = c
					pix[offset+3] = 255
				}
			}
		}
	}

	return img
}

package qr

import (
	"aboutblank/qr-code/bitreader"
	"image"
	"math"
)

type Version uint8

type QRCode struct {
	Version      Version
	EcLevel      ErrorCorrectionLevel
	ModuleMatrix [][]Module
	Mask         int

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

var formatInfo = [4][8]uint16{
	// EC_L
	{
		0b111011111000100, // pattern 0
		0b111001011110011, // pattern 1
		0b111110110101010,
		0b111100010011101,
		0b110011000101111,
		0b110001100011000,
		0b110110001000001,
		0b110100101110110,
	},
	// EC_M
	{
		0b101010000010010,
		0b101000100100101,
		0b101111001111100,
		0b101101101001011,
		0b100010111111001,
		0b100000011001110,
		0b100111110010111,
		0b100101010100000,
	},
	// EC_Q
	{
		0b011010101011111,
		0b011000001101000,
		0b011111100110001,
		0b011101000000110,
		0b010010010110100,
		0b010000110000011,
		0b010111011011010,
		0b010101111101101,
	},
	// EC_H
	{
		0b001011010001001,
		0b001001110111110,
		0b001110011100111,
		0b001100111010000,
		0b000011101100010,
		0b000001001010101,
		0b000110100001100,
		0b000100000111011,
	},
}

var versionInfo = map[int]uint32{
	7:  0b000111110010010100,
	8:  0b001000010110111100,
	9:  0b001001101010011001,
	10: 0b001010010011010011,
	11: 0b001011101111110110,
	12: 0b001100011101100010,
	13: 0b001101100001000111,
	14: 0b001110011000001101,
	15: 0b001111100100101000,
	16: 0b010000101101111000,
	17: 0b010001010001011101,
	18: 0b010010101000010111,
	19: 0b010011010100110010,
	20: 0b010100100110100110,
	21: 0b010101011010000011,
	22: 0b010110100011001001,
	23: 0b010111011111101100,
	24: 0b011000111011000100,
	25: 0b011001000111100001,
	26: 0b011010111110101011,
	27: 0b011011000010001110,
	28: 0b011100110000011010,
	29: 0b011101001100111111,
	30: 0b011110110101110101,
	31: 0b011111001001010000,
	32: 0b100000100111010101,
	33: 0b100001011011110000,
	34: 0b100010100010111010,
	35: 0b100011011110011111,
	36: 0b100100101100001011,
	37: 0b100101010000101110,
	38: 0b100110101001100100,
	39: 0b100111010101000001,
	40: 0b101000110001101001,
}

func New(version Version, ecLevel ErrorCorrectionLevel) *QRCode {
	qr := &QRCode{}
	qr.Version = version
	qr.EcLevel = ecLevel
	qr.Mask = -1

	size := int(21 + 4*(version-1))
	qr.size = size

	qr.ModuleMatrix = make([][]Module, qr.size)
	for i := range qr.ModuleMatrix {
		qr.ModuleMatrix[i] = make([]Module, qr.size)
	}

	qr.formatPositions = [30][2]int{
		// Top-left area
		{8, 0}, {8, 1}, {8, 2}, {8, 3}, {8, 4}, {8, 5}, {8, 7},
		{8, 8}, {7, 8}, {5, 8}, {4, 8}, {3, 8}, {2, 8}, {1, 8}, {0, 8},

		// Mirror area
		{size - 1, 8}, {size - 2, 8}, {size - 3, 8}, {size - 4, 8}, {size - 5, 8}, {size - 6, 8}, {size - 7, 8}, {size - 8, 8},
		{8, size - 7}, {8, size - 6}, {8, size - 5}, {8, size - 4}, {8, size - 3}, {8, size - 2}, {8, size - 1},
	}

	qr.versionPositions = [36][2]int{
		// Top-right block (3x6)
		{size - 11, 0}, {size - 10, 0}, {size - 9, 0},
		{size - 11, 1}, {size - 10, 1}, {size - 9, 1},
		{size - 11, 2}, {size - 10, 2}, {size - 9, 2},
		{size - 11, 3}, {size - 10, 3}, {size - 9, 3},
		{size - 11, 4}, {size - 10, 4}, {size - 9, 4},
		{size - 11, 5}, {size - 10, 5}, {size - 9, 5},

		// Bottom-left block (6x3)
		{0, size - 11}, {1, size - 11}, {2, size - 11}, {3, size - 11}, {4, size - 11}, {5, size - 11},
		{0, size - 10}, {1, size - 10}, {2, size - 10}, {3, size - 10}, {4, size - 10}, {5, size - 10},
		{0, size - 9}, {1, size - 9}, {2, size - 9}, {3, size - 9}, {4, size - 9}, {5, size - 9},
	}

	return qr
}

func (qr *QRCode) Clone() *QRCode {
	size := len(qr.ModuleMatrix)

	newMatrix := make([][]Module, size)
	for i := range newMatrix {
		newMatrix[i] = make([]Module, size)
		copy(newMatrix[i], qr.ModuleMatrix[i])
	}

	return &QRCode{
		Version:      qr.Version,
		ModuleMatrix: newMatrix,
	}
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
		if !reader.HasData() {
			break
		}

		val := ValueWhite
		if reader.Pop() {
			val = ValueBlack
		}

		qr.setModule(pos[0], pos[1], val, false)
	}
}

func (qr *QRCode) WriteFormatInfo() {
	info := formatInfo[qr.EcLevel][qr.Mask]

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

func (qr *QRCode) ApplyBestMask() {
	bestScore := math.MaxInt
	bestMask := 0

	for mask := range 8 {
		clone := qr.Clone()
		clone.ApplyMask(mask)
		score := clone.Score()

		if score < bestScore {
			bestScore = score
			bestMask = mask
		}
	}

	qr.ApplyMask(bestMask)
}

func (qr *QRCode) ApplyMask(mask int) {
	qr.Mask = mask

	size := len(qr.ModuleMatrix)
	for y := range size {
		for x := range size {
			mod := &qr.ModuleMatrix[y][x]

			if mod.Reserved {
				continue
			}

			if maskApplies(mask, x, y) {
				switch mod.Value {
				case ValueBlack:
					mod.Value = ValueWhite
				case ValueWhite:
					mod.Value = ValueBlack
				}
			}
		}
	}
}

func (qr *QRCode) Score() int {
	size := len(qr.ModuleMatrix)
	score := 0

	// Rule 1: rows
	for y := range size {
		run := 1
		for x := 1; x < size; x++ {
			if qr.ModuleMatrix[y][x].Value == qr.ModuleMatrix[y][x-1].Value {
				run++
			} else {
				if run >= 5 {
					score += 3 + (run - 5)
				}
				run = 1
			}
		}
		if run >= 5 {
			score += 3 + (run - 5)
		}
	}

	// Rule 1: columns
	for x := range size {
		run := 1
		for y := 1; y < size; y++ {
			if qr.ModuleMatrix[y][x].Value == qr.ModuleMatrix[y-1][x].Value {
				run++
			} else {
				if run >= 5 {
					score += 3 + (run - 5)
				}
				run = 1
			}
		}
		if run >= 5 {
			score += 3 + (run - 5)
		}
	}

	// Rule 2: 2x2 blocks
	for y := 0; y < size-1; y++ {
		for x := 0; x < size-1; x++ {
			v := qr.ModuleMatrix[y][x].Value
			if qr.ModuleMatrix[y+1][x].Value == v &&
				qr.ModuleMatrix[y][x+1].Value == v &&
				qr.ModuleMatrix[y+1][x+1].Value == v {
				score += 3
			}
		}
	}

	// Rule 4: dark ratio
	dark := 0
	total := size * size
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if qr.ModuleMatrix[y][x].Value == ValueBlack {
				dark++
			}
		}
	}
	percent := (dark * 100) / total
	diff := percent - 50
	if diff < 0 {
		diff = -diff
	}
	score += (diff / 5) * 10

	return score
}

func maskApplies(mask, x, y int) bool {
	switch mask {
	case 0:
		return (x+y)%2 == 0
	case 1:
		return y%2 == 0
	case 2:
		return x%3 == 0
	case 3:
		return (x+y)%3 == 0
	case 4:
		return (y/2+x/3)%2 == 0
	case 5:
		return (x*y)%2+(x*y)%3 == 0
	case 6:
		return ((x*y)%2+(x*y)%3)%2 == 0
	case 7:
		return ((x+y)%2+(x*y)%3)%2 == 0
	}
	return false
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

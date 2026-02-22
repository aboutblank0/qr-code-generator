package qr

import "math"

func (qr *QRCode) ApplyBestMask() {
	bestScore := math.MaxInt
	bestMask := 0

	for mask := range 8 {
		clone := qr.Clone()
		clone.ApplyMask(mask)
		score := clone.ScoreMask()

		if score < bestScore {
			bestScore = score
			bestMask = mask
		}
	}

	qr.ApplyMask(bestMask)
}

func (qr *QRCode) ApplyMask(mask int) {
	qr.mask = mask

	size := len(qr.moduleMatrix)
	for y := range size {
		for x := range size {
			mod := &qr.moduleMatrix[y][x]

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

func (qr *QRCode) ScoreMask() int {
	size := len(qr.moduleMatrix)
	score := 0

	// Rule 1: rows
	for y := range size {
		run := 1
		for x := 1; x < size; x++ {
			if qr.moduleMatrix[y][x].Value == qr.moduleMatrix[y][x-1].Value {
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
			if qr.moduleMatrix[y][x].Value == qr.moduleMatrix[y-1][x].Value {
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
			v := qr.moduleMatrix[y][x].Value
			if qr.moduleMatrix[y+1][x].Value == v &&
				qr.moduleMatrix[y][x+1].Value == v &&
				qr.moduleMatrix[y+1][x+1].Value == v {
				score += 3
			}
		}
	}

	// Rule 4: dark ratio
	dark := 0
	total := size * size
	for y := range size {
		for x := range size {
			if qr.moduleMatrix[y][x].Value == ValueBlack {
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

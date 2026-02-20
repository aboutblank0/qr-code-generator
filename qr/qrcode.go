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

func (qr *QRCode) GenerateImage(scale int) *image.RGBA {
	size := qr.GetSize()
	w, h := size*scale, size*scale
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	pix := img.Pix
	stride := img.Stride

	for x := range size {
		for y := range size {
			c := byte(255) // white
			if qr.ModuleMatrix[x][y].Value {
				c = 0
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

func (qr *QRCode) AddFinderPatterns() {
	size := qr.GetSize()

	// Top left
	qr.addSquare(0, 0, 7, false)
	qr.addSquare(0 + 2, 0 + 2, 3, true)

	// Bottom Left
	qr.addSquare(0, size-7, 7, false)
	qr.addSquare(0 + 2, size-7 + 2, 3, true)

	// Top Right
	qr.addSquare(size-7, 0, 7, false)
	qr.addSquare(size-7 + 2, 0 + 2, 3, true)

}

func (qr *QRCode) addSquare(startX, startY, size int, fill bool) {
	for y := range size {
		for x := range size {
			if fill {
				qr.ModuleMatrix[startX+x][startY+y].Value = true
				continue
			}

			if x == 0 || x == size-1 || y == 0 || y == size-1 {
				qr.ModuleMatrix[startX+x][startY+y].Value = true
			}
		}
	}
}


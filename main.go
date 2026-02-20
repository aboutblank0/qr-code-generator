package main

import (
	"aboutblank/qr-code/qr"
	"fmt"
	"image"
	"image/png"
	"os"
)

func main() {
	//DELETE ME
	finalMessage := qr.GenerateQRCode("HELLO WORLD", qr.Alphanumeric, qr.EC_M)
	fmt.Printf("%d\n", finalMessage)

	//res := qr.PolyMultiply([]byte{1, 0, 0, 1}, []byte{1, 0, 0, 1})
	//fmt.Printf("%d\n", res)
	//mult := gf256.Multiply(0b00000001, 0b00000001)
	//fmt.Printf("%08b\n", mult)

	qrCode := qr.New(qr.Version(1))

	qrCode.AddFinderPatterns()
	image := qrCode.GenerateImage(4)
	SaveImage(image, "qrcode.png")
}

func SaveImage(image *image.RGBA, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = png.Encode(f, image)
	if err != nil {
		return nil
	}
	return nil
}

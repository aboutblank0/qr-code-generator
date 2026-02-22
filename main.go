package main

import (
	"aboutblank/qr-code/qr"
	"image"
	"image/png"
	"os"
)

func main() {
	qrCode := qr.GenerateQRCode("1234567890", qr.Encode_Numeric, qr.EC_High)
	image := qrCode.GenerateImage(10)
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

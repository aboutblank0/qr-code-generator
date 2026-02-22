package main

import (
	"aboutblank/qr-code/qr"
	"image"
	"image/png"
	"os"
)

func main() {
	data := qr.GenerateQRCode("HELLO WORLD", qr.Alphanumeric, qr.EC_Q)
	qrCode := qr.New(qr.Version(1), qr.EC_Q)
	qrCode.Test(data)

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

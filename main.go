package main

import (
	"aboutblank/qr-code/qr"
	"fmt"
	"image"
	"image/png"
	"os"
)

func main() {
	GenerateEveryErrorLevel("Hi!")
}

func GenerateEveryErrorLevel(data string) {
	qrCode := qr.GenerateQRCode(data, qr.EC_Low)
	image := qrCode.GenerateImage(10)
	SaveImage(image, "qrcode_low.png")
	fmt.Println()

	qrCode = qr.GenerateQRCode(data, qr.EC_Medium)
	image = qrCode.GenerateImage(10)
	SaveImage(image, "qrcode_medium.png")
	fmt.Println()

	qrCode = qr.GenerateQRCode(data, qr.EC_Quartile)
	image = qrCode.GenerateImage(10)
	SaveImage(image, "qrcode_quartile.png")
	fmt.Println()

	qrCode = qr.GenerateQRCode(data, qr.EC_High)
	image = qrCode.GenerateImage(10)
	SaveImage(image, "qrcode_high.png")
	fmt.Println()
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

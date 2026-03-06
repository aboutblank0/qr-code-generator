package main

import (
	"aboutblank/qr-code/qr"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
)

func main() {
	var helpFlag = flag.Bool("help", false, "Display help information")
	var scaleFlag = flag.Int("scale", 10, "Scale factor for the generated QR code image")
	var outputFlag = flag.String("output", "qrcode.png", "Output file name for the generated QR code image")
	var versionOverrideFlag = flag.Int("version", 0, "Override QR code version (1-40)")

	flag.Parse()

	// Show help
	if *helpFlag {
		PrintHelp()
		return
	}

	// Check there is content to write
	if flag.NArg() < 1 {
		fmt.Println("ERR: No content provided.")
		fmt.Println("Usage: qrgen [options] <content>")
		return
	}

	if *versionOverrideFlag < 0 || *versionOverrideFlag > 40 {
		fmt.Println("ERR: Version must be between 1 and 40.")
		return
	}

	content := flag.Arg(0)
	qrCode := qr.GenerateQRCode(content, qr.EC_Medium, *versionOverrideFlag)
	image := qrCode.GenerateImage(*scaleFlag)
	err := SaveImage(image, *outputFlag)

	if err != nil {
		fmt.Println("ERR: Failed to save image:", err)
		return
	}
}

func PrintHelp() {
	fmt.Println("Usage: qrgen [options] <content>")
	fmt.Println("Options:")
	flag.PrintDefaults()
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
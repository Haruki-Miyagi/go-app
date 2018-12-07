package main

import (
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func main() {
	qrCode, _ := qr.Encode("https://github.com/haruki-miyagi", qr.M, qr.Auto)

	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	file, _ := os.Create("mygithub.png")
	defer file.Close()

	png.Encode(file, qrCode)
}

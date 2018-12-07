package main

import (
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func main() {
	// 元となるimageファイル
	file, err := os.Open("./images/image01.jpg")
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// 引数(幅, 高さ, 画像data, アルゴリズムの設定)
	m := resize.Resize(1000, 0, img, resize.Lanczos3)

	out, err := os.Create("./resized_images/r_image01.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	jpeg.Encode(out, m, nil)
}

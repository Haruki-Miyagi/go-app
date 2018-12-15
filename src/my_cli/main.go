package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "my_cli"
	app.Usage = "my command"

	app.Commands = []cli.Command{
		{
			// qrcode
			Name:    "qrcode",
			Aliases: []string{"q"},
			Usage:   "make qrcode",
			Action: func(c *cli.Context) error {
				// cにはqrcodeの後の文字列を取得する
				fmt.Println("qrcode complete: ", c.Args().First())
				qrCode, _ := qr.Encode(c.Args().First(), qr.M, qr.Auto)
				qrCode, _ = barcode.Scale(qrCode, 200, 200)
				file, _ := os.Create("qrcode.png")
				defer file.Close()
				png.Encode(file, qrCode)
				return nil
			},
		},
		{
			// githubのリポジトリ名の取得自作
			// 使用例) my_cli c https://github.com/Haruki-Miyagi?tab=repositories
			Name:    "crawler",
			Aliases: []string{"c"},
			Usage:   "crawler repositories name of github",
			Action: func(c *cli.Context) error {
				res, err := http.Get(c.Args().First())
				if err != nil {
					log.Fatal(err)
				}
				defer res.Body.Close()
				if res.StatusCode != 200 {
					log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
				}
				fmt.Println("repositories name complete URL: ", c.Args().First())
				// Load the HTML document
				doc, err := goquery.NewDocumentFromReader(res.Body)
				if err != nil {
					log.Fatal(err)
				}

				// Find the review items
				doc.Find(".d-inline-block .mb-1").Each(func(i int, s *goquery.Selection) {
					// For each item found, get the band and title
					band := s.Find("a")
					link, _ := band.Last().Attr("href")

					fmt.Printf("All repositories name URL : %s \n", link)
				})
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

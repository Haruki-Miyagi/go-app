package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
)

func main() {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	res, err := http.Get("https://github.com/Haruki-Miyagi?tab=repositories")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

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

		// xlsxファイルを作成
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = link
		err = file.Save("MyXLSXFile.xlsx")
		if err != nil {
			fmt.Printf(err.Error())
		}

		fmt.Printf("aタグ %s \n", link)
	})
}

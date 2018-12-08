package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
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
		fmt.Printf("aタグ %s \n", link)
	})
}

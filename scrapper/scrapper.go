package scrapper

import (
	"fmt"
	"github.com/gocolly/colly"
)
type ConfigHttp struct {
	URL        string
	Timeout    int
	Delay      float64
	Threads    int
	OutputType int // 6: json, 2: csv, 4: xml, 8: text

}

var DefaultConfig ConfigHttp  =   ConfigHttp{
	URL: "",
	Timeout: 10,
	Delay: 0.5,
	Threads: 1,
	OutputType: 2,
}


func VisitUrl(url string) {
	if url == "" {
		fmt.Println("URL is required")
		return
	}
	c := colly.NewCollector()
	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", h.Text, link)
	})
	c.Visit(url)
}

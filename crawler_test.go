package cuvva

import (
	"net/url"
	"testing"
)

func TestCrawl(t *testing.T) {
	URL, err := url.Parse("https://www.cuvva.com")
	if err != nil {
		return
	}
	crawler := NewCrawler(URL)
	crawler.Crawl()
	crawler.Print()
}

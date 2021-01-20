package cuvva

import (
	"net/url"
	"testing"
)

func TestCrawl(t *testing.T) {
	URL, _ := url.Parse("https://www.cuvva.com")
	crawler := NewCrawler()
	crawler.Crawl(URL)
	crawler.Print()
}

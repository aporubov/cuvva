package cuvva

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// Resource struct represents any web resource
type Resource struct {
	URL          *url.URL
	References   []*Resource
	ContentType  string
	StatusCode   int
	LastModified string
	Error        error
}

// Crawler struct represents main crawling context
type Crawler struct {
	httpClient *http.Client
	URL        *url.URL
	// Entry point to the graph of collected resources
	entry     *Resource
	resources map[string]*Resource
	mux       sync.Mutex
}

// NewCrawler instantiates a new Crawler
func NewCrawler(URL *url.URL) *Crawler {
	return &Crawler{
		httpClient: &http.Client{Timeout: time.Minute},
		URL:        URL,
		resources:  make(map[string]*Resource),
	}
}

// Crawl performs the crawling
func (c *Crawler) Crawl() {
	c.resources = make(map[string]*Resource)
	c.entry = c.process(c.URL)
}

// Print prints generated sitemap
func (c *Crawler) Print() {
	for _, resource := range c.resources {
		if resource.Error != nil || !resource.isWebPage() {
			continue
		}
		fmt.Printf("(webpage) %v\n", resource.URL.String())
		fmt.Printf("    Links\n")
		for _, ref := range resource.References {
			if ref.Error != nil || !ref.isWebPage() {
				continue
			}
			fmt.Printf("        (webpage) %v\n", ref.URL.String())
		}
		fmt.Printf("    Assets\n")
		for _, ref := range resource.References {
			if ref.Error != nil || ref.isWebPage() {
				continue
			}
			fmt.Printf("        (%v) %v\n", ref.ContentType, ref.URL.String())
		}
	}
}

func (c *Crawler) process(URL *url.URL) *Resource {
	if !c.shouldVisit(URL) {
		// Ignore any external resources
		return nil
	}

	resource, isVisited := c.isVisited(URL)
	if isVisited {
		return resource
	}

	resp, err := c.fetch(URL)
	if err != nil {
		fmt.Println(err)
		// Keep a failed to fetch resource in the registry
		// of known resources so it's not attempted to re-crawl
		// Keep occurred HTTP error on a resource to indicate
		// the resource wasn't fetched successfully
		// TODO: employ retry or other HTTP error handling
		// strategy here if needed
		resource.Error = err
		return resource
	}

	refURLs := c.analyze(resp, resource)
	ch := make(chan *Resource)

	for _, refURL := range refURLs {
		go func(u *url.URL) {
			ch <- c.process(u)
		}(refURL)
	}

	for i := 0; i < len(refURLs); i++ {
		refResource := <-ch
		if refResource != nil {
			resource.References = append(resource.References, refResource)
		}
	}

	fmt.Printf("processed %v\n", resource.URL.String())

	return resource
}

func (c *Crawler) shouldVisit(URL *url.URL) bool {
	return URL.Hostname() == c.URL.Hostname()
}

func (c *Crawler) isVisited(URL *url.URL) (*Resource, bool) {
	// Obtain exclusive access to the registry of known resources
	c.mux.Lock()
	defer c.mux.Unlock()
	existing, isVisited := c.resources[URL.String()]
	if isVisited {
		return existing, true
	}
	new := &Resource{URL: URL}
	c.resources[URL.String()] = new
	return new, false
}

func (c *Crawler) fetch(URL *url.URL) (*http.Response, error) {
	return c.httpClient.Get(URL.String())
}

func (c *Crawler) analyze(resp *http.Response, resource *Resource) (refURLs []*url.URL) {
	defer resp.Body.Close()
	resource.ContentType = resp.Header.Get("Content-Type")
	resource.LastModified = resp.Header.Get("Last-Modified")
	resource.StatusCode = resp.StatusCode
	if !resource.isWebPage() {
		return
	}
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			if tokenizer.Err() != io.EOF {
				// Ignore this for now
				fmt.Println("Error parsing HTML document:", tokenizer.Err())
			}
			return unique(refURLs)
		}
		token := tokenizer.Token()
		switch tokenType {
		case html.StartTagToken, html.SelfClosingTagToken:
			refURL, err := c.extract(token)
			if err != nil {
				continue
			}
			refURLs = append(refURLs, refURL)
		}
	}
}

func (c *Crawler) extract(token html.Token) (*url.URL, error) {
	for _, attr := range token.Attr {
		switch attr.Key {
		case "href", "src":
			trimmed := strings.TrimSpace(attr.Val)
			URL, err := url.Parse(trimmed)
			if err != nil {
				return nil, err
			}
			if !URL.IsAbs() {
				URL = c.URL.ResolveReference(URL)
			}
			return URL, nil
		}
	}
	return nil, errors.New("URL not found")
}

func (r *Resource) isWebPage() bool {
	return strings.Contains(r.ContentType, "text/html")
}

func unique(s []*url.URL) []*url.URL {
	seen := map[string]bool{}
	result := []*url.URL{}
	for _, v := range s {
		if seen[v.String()] != true {
			seen[v.String()] = true
			result = append(result, v)
		}
	}
	return result
}

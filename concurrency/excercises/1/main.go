package main

/*
	Aufgabenstellung:

*/

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
)

type Fetcher interface {
	Fetch(url string) ([]byte, int, error)
}

type HttpFetcher struct{}

type Result struct {
	URL        string
	StatusCode int
	Length     int
	Err        error
}

func main() {
	fetcher := HttpFetcher{}

	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/uuid",
		"https://httpbin.alex/penis",
	}

	results := FetchAllURLs(fetcher, urls)
	for _, res := range results {
		sort.Slice(results, func(i, j int) bool {
			return results[i].URL < results[j].URL
		})
		fmt.Println(res.URL, res.StatusCode, res.Length, res.Err)
	}
}

// hiermit ist die Bedingung erfüllt
func (h HttpFetcher) Fetch(url string) ([]byte, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // err wird neu belegt, aber Get-Fehler ist schon behandelt
	if err != nil {
		return nil, 0, fmt.Errorf("read body failed: %w", err)
	}

	return body, resp.StatusCode, nil
}

func FetchAll(fetcher Fetcher, urls []string) []Result {
	var wg sync.WaitGroup
	// deutlich einfacher scheinbar:
	//channel := make(chan Result, len(urls))

	results := []Result{}
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			body, status, err := fetcher.Fetch(u)
			res := Result{u, status, len(body), err}
			//channel <- res
			results = append(results, res)
		}(url)
	}
	wg.Wait()
	return results
}

func FetchAllURLs(fetcher Fetcher, urls []string) []Result {
	var wg sync.WaitGroup
	channel := make(chan Result, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			body, status, err := fetcher.Fetch(u)
			channel <- Result{u, status, len(body), err}
		}(url)
	}

	wg.Wait()
	close(channel)

	// somit werden dann eben race condition warnings unterbunden
	var results []Result
	for res := range channel {
		results = append(results, res)
	}
	return results
}

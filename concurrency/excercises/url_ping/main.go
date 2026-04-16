package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	//"golang.org/x/tools/go/analysis/passes/waitgroup"
)

type Result struct {
	URL    string
	Status int
	Err    error
}

func checkAlex(urls []string) []Result {
	// TODO: jede URL in einer eigenen Goroutine prüfen
	// HTTP GET, Status-Code und Fehler zurückgeben
	// Ergebnisse sammeln ohne Race Condition
	// mein Ansatz nicht so toll aber auch full lost... output garnicht so schlecht

	wg := sync.WaitGroup{}
	ch := make(chan Result, len(urls))
	var res []Result
	for _, url := range urls {
		var result Result
		wg.Add(1)
		response, err := http.Get(url)
		if err != nil {
			result.Err = err
		} else {
			result.Status = response.StatusCode
			result.Err = nil
			result.URL = url
		}
		defer response.Body.Close()
		go func(url string) {
			defer wg.Done()
			ch <- result
		}(url)

	}

	wg.Wait()
	close(ch)

	for r := range ch {
		res = append(res, r)
	}
	//print("finished")

	return res
}

func checkAll(urls []string) []Result {
	ch := make(chan Result, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			result := checkOne(url)
			ch <- result
		}(url) // muss vom loop sein sonst Fehler
	}

	wg.Wait()
	close(ch)

	var res []Result
	for r := range ch {
		res = append(res, r)
	}
	return res
}

func checkOne(url string) Result {
	ch := make(chan Result, 1) // buffered, oder fixed length

	go func() {
		result := Result{URL: url}

		resp, err := http.Get(url)
		if err != nil {
			result.Err = err
			ch <- result
			return
		}
		result.Status = resp.StatusCode
		if resp.StatusCode >= 500 {
			result.Err = fmt.Errorf("server error: %d", resp.StatusCode)
		}
		ch <- result // ← das hat gefehlt
		defer resp.Body.Close()
	}()

	select {
	case result := <-ch:
		return result
	case <-time.After(2 * time.Second):
		return Result{URL: url, Status: http.StatusRequestTimeout, Err: fmt.Errorf("timeout after 2s")}
	}
}

func main() {
	urls := []string{
		"https://httpbin.org/delay/5",
		"https://example.com",
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/500",
	}
	results := checkAll(urls)
	for _, r := range results {
		fmt.Printf("%s → %d %v\n", r.URL, r.Status, r.Err)
	}
}

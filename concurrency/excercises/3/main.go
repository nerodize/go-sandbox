package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	URL        string // hier vllt noch das json Zeug
	StatusCode int
	Duration   time.Duration
	Error      error
}

func main() {
	urls := []string{
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/3",
		"https://httpbin.org/delay/5",
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/500",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // sieht fast nach gängigem pattern aus

	resultsCh := CheckURLs(ctx, urls)

	for i := 0; i < len(urls); i++ {
		select {
		case result := <-resultsCh:
			if result.Error != nil {
				fmt.Printf("✗ %-45s %s\n", result.URL, result.Error) // what the...
			} else {
				fmt.Printf("✓ %-45s %d  %s\n", result.URL, result.StatusCode, result.Duration)
			}
		case <-ctx.Done():
			fmt.Println("globaler timeout - abbruch")
			return
		}
	}
}

func CheckURLs(ctx context.Context, urls []string) chan Response {
	// hier müsste ja dann noch etwas mit dem respChan passieren sonst ist er ja leer
	// höchstens wenn er in der main also im Param befüllt wird.
	responseChan := make(chan Response, len(urls))
	client := http.Client{
		Timeout: 1500 * time.Millisecond,
	}

	for _, url := range urls {
		go func(url string) {
			resultChan := make(chan Response, 1) // buffered, pro goroutine ein Channel

			go func() {
				start := time.Now()
				req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
				if err != nil {
					resultChan <- Response{URL: url, Error: err}
					return
				}

				resp, err := client.Do(req)
				duration := time.Since(start)
				if err != nil {
					resultChan <- Response{URL: url, Duration: duration, Error: err}
					return
				}
				defer resp.Body.Close()

				resultChan <- Response{URL: url, StatusCode: resp.StatusCode, Duration: duration, Error: err}

			}()

			select {
			case result := <-resultChan:
				responseChan <- result
			case <-ctx.Done():
				responseChan <- Response{URL: url, Error: fmt.Errorf("global timeout after: %w", ctx.Err())}
			}
		}(url)

	}

	// und das macht dann ebenso wenig sinn.
	return responseChan
}

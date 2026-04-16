package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	url := "https://jsonplaceholder.typicode.com/posts"
	channel, errChannel1 := FetchJSON(url)
	// bei falsche url: erst hier der fehler dann...
	postChannel, errChannel2 := ParseJson(channel)
	filteredChannel, errChannel3 := Filter(postChannel)

	for post := range filteredChannel {
		fmt.Println(post.UserID, post.Title)
	}

	// Ausgabe des Error slices...
	for _, errCh := range []<-chan error{errChannel1, errChannel2, errChannel3} {
		if err := <-errCh; err != nil {
			fmt.Println(err)
		}
	}

}

// stage 3
func Filter(in <-chan Post) (<-chan Post, <-chan error) {
	out := make(chan Post)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)

		for post := range in {
			if post.UserID == 1 {
				out <- post
			}
		}

	}()
	return out, errCh
}

// stage 2
func ParseJson(in <-chan []byte) (<-chan Post, <-chan error) {
	out := make(chan Post) // oder iwie mit len(posts))
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)
		jsonValue := <-in

		var posts []Post
		err := json.Unmarshal(jsonValue, &posts)
		if err != nil {
			errCh <- fmt.Errorf("parsing: %w", err)
			return
		}

		for _, post := range posts {
			out <- post
		}
	}()

	return out, errCh
}

// stage 1
func FetchJSON(url string) (<-chan []byte, <-chan error) {
	channel := make(chan []byte, 1) // könnte problematisch werden wegen buffered bzw. !buffered
	errCh := make(chan error, 1)
	go func() {
		defer close(channel)
		defer close(errCh)

		resp, err := http.Get(url)
		if err != nil {
			errCh <- fmt.Errorf("Fehler beim fetchen: %w", err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body) // err wird neu belegt, aber Get-Fehler ist schon behandelt
		if err != nil {
			errCh <- fmt.Errorf("Fehler: %w", err)
			return
		}

		channel <- body
	}()
	return channel, errCh
}

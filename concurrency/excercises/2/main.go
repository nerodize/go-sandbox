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
	channel := FetchJSON(url)

	filterChannel, errChannel := ParseJson(channel)

	print(filterChannel, errChannel)
	/*
		jsonValue := <-channel
		var posts []Post
		err := json.Unmarshal(jsonValue, &posts)
		if err != nil {
			fmt.Errorf("Fehler beim Parsen: %w", err)
		}

		postCh := make(chan Post)
		for _, post := range posts {
			postCh <- post
		}
	*/

	// hier die funktion mit posts als param

}

func ParseJson(in <-chan []byte) (<-chan Post, <-chan error) {
	out := make(chan Post) // oder iwie mit len(posts))
	errCh := make(chan error)

	go func() {
		defer close(out)
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
func FetchJSON(url string) <-chan []byte {
	channel := make(chan []byte, 1) // könnte problematisch werden wegen buffered bzw. !buffered

	go func() {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Errorf("Fehler beim fetchen: %w", err)
		}

		body, err := io.ReadAll(resp.Body) // err wird neu belegt, aber Get-Fehler ist schon behandelt
		if err != nil {
			fmt.Errorf("Fehler: %w", err)
		}

		channel <- body
		close(channel)
	}()
	return channel
}

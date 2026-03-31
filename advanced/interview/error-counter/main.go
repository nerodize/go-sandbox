package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	file, err := os.Open("C:/dev/learning/go/sandbox/advanced/interview/error-counter/app.log")

	if err != nil {
		fmt.Println("Error opening file: ", err)
	}
	defer file.Close()

	errorCounts := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "ERROR") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}

			timestampStr := parts[0] + " " + parts[1]

			layout := "2006-01-02 15:04:05"
			t, err := time.Parse(layout, timestampStr)
			if err != nil {
				continue
			}
			// format
			minuteStr := t.Format("2006-01-02 15:04")
			// hier wird die map befüllt... ["Minute": "value"(als Inkrement)]
			errorCounts[minuteStr]++
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file", err)
		return
	}

	type MinuteCount struct {
		Minute time.Time // probably smarter to call this minute...
		Count  int
	}

	var counts []MinuteCount
	for k, v := range errorCounts {
		t, err := time.Parse("2006-01-02 15:04", k)
		if err != nil {
			fmt.Println("error while parsing string to time", err)
			continue
		}

		counts = append(counts, MinuteCount{Minute: t, Count: v})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Minute.Before(counts[j].Minute)
	})

	for _, minuteCount := range counts {
		fmt.Printf("%s -> %d\n", minuteCount.Minute.Format("2006-01-02 15:04"), minuteCount.Count)
	}
}

package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	// slices => dynamisch und arrays fest
	output := TwoSums([]int{2, 7, 11, 15}, 18)
	fmt.Println(output)

	text := "Go is great and Go is fast and fast is good"
	res := CountWords(text)
	fmt.Printf("kek: %v", res)
	fmt.Print("\n------\n")

	wordCount := map[string]int{
		"go":    5,
		"is":    3,
		"great": 1,
		"fast":  4,
	}

	top3 := TopNWords(wordCount, 3)
	for _, line := range top3 {
		fmt.Println(line)
	}

	// nur das erste Element printen
	nums := []int{1, 3, 4, 5}
	print(nums[:1])

}

func MergeIntervals(intervals [][]int) [][]int {

	if len(intervals) == 0 {
		return intervals
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	merged := [][]int{}
	current := intervals[0]

	for i := 1; i < len(intervals); i++ { // Start bei 1, nicht 0!
		interval := intervals[i]

		// >= sonst kein Überlappen
		if current[1] >= interval[0] {
			current[1] = interval[1]
		} else {
			merged = append(merged, current)
			current = interval
		}

	}
	merged = append(merged, current)
	return merged
}

func TwoSums(nums []int, target int) []int {
	seen := make(map[int]int)
	// scheinbar nötig teilweise beim
	defer func() {
		fmt.Printf("map contents: %v\n", seen)
	}()

	for i, num := range nums {
		complement := target - num
		if index, exists := seen[complement]; exists {
			return []int{index, i}
		} else {
			seen[num] = i
		}
	}
	return nil
}

func TwoSumsDumb(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{nums[i], nums[j]}
			}
		}
	}
	return nil
}

/*
func CountWords(text string) map[string]int {
	returnMap := make(map[string]int)
	seen := make(map[string]int)

	tt := strings.ToLower(text)
	wordList := strings.Fields(tt)
	for _, word := range wordList {
		fmt.Println(word)
		wordCount := 1
		if _, exists := seen[word]; exists {
			wordCount++
			returnMap[word] = wordCount
		}
	}

	return returnMap
}
*/

func CountWords(input string) map[string]int {
	returnMap := make(map[string]int)

	wordList := strings.Fields(strings.ToLower(input)) // warum überhaupt so vorbereiten?

	for _, word := range wordList {
		returnMap[word]++
	}

	return returnMap
}

func CountWordss(input string) map[string]int {
	returnMap := make(map[string]int)

	wordList := strings.Fields(strings.ToLower(input)) // warum überhaupt so vorbereiten?

	for _, word := range wordList {
		if count, exists := returnMap[word]; exists {
			returnMap[word] = count + 1
		} else {
			returnMap[word] = 1
		}

	}

	return returnMap
}

// this one with a sctruct but only as an extra
func TopNWords(wordCount map[string]int, n int) []string {
	returnArray := []string{}
	if len(wordCount) < n {
		n = len(wordCount)
	}

	type WordCount struct {
		Word  string
		Count int
	}

	pairs := []WordCount{}
	for word, count := range wordCount {
		pairs = append(pairs, WordCount{word, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Count > pairs[j].Count
	})

	for i := 0; i < n && i < len(pairs); i++ {
		formatted := fmt.Sprintf("%s: %d", pairs[i].Word, pairs[i].Count)
		returnArray = append(returnArray, formatted)
	}

	return returnArray
}

func TopNWordsNoCount(wordCount map[string]int, n int) []string {
	//returnArray := []string{}

	if len(wordCount) < n {
		n = len(wordCount)
	}

	words := []string{}
	for word := range wordCount {
		words = append(words, word)
	}

	sort.Slice(words, func(i, j int) bool {
		return wordCount[words[i]] > wordCount[words[j]]
	})

	return words[:n] // die ersten n Einträge zurückgeben
}

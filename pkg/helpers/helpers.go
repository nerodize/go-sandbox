package helpers

import (
	"fmt"
	"strings"
	"unicode"
)

func HasUppercase(s string) bool {
	for _, char := range s {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

func HasUppercaseSlow(s string) bool {
	hasUpper := false
	for _, char := range s {
		// sonst kein missmatch
		if string(char) != strings.ToLower(string(char)) {
			hasUpper = true
			fmt.Printf("returns true on this letter: %c", char)
			break
		}
	}
	return hasUpper
}

func HasDigit(s string) bool {
	for _, char := range s {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

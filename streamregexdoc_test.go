package streamregex

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

func ExampleRegex_FindReader() {
	// Create string
	data := `0123456789this is a stream    of data with lots of trailing information`
	stream := strings.NewReader(data)

	// Build regex
	regexInt := regexp.MustCompile(`stream\s+of`)
	regex := NewRegex(regexInt)

	// Find matches
	matchedData := regex.FindReader(context.Background(), stream)
	matches := 0
	for match := range matchedData {
		matches++
		fmt.Println(string(match))
	}

	// Output: stream    of
}

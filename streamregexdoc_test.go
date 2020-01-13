package streamregex

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

func ExampleFindReader() {
	// Create string
	data := `0123456789this is a stream    of data with lots of trailing information`
	stream := strings.NewReader(data)

	// Build regex
	regex := regexp.MustCompile(`stream\s+of`)

	// Find matches
	matchedData := FindReader(context.Background(), regex, 100, stream)
	matches := 0
	for match := range matchedData {
		matches++
		fmt.Println(match)
	}

	// Output: stream    of
}

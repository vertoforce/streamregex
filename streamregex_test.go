package streamregex

import (
	"context"
	"regexp"
	"strings"
	"testing"
)

func TestFindReader(t *testing.T) {
	// Create string
	data := `0123456789this is a stream    of data with lots of trailing information`
	stream := strings.NewReader(data)

	// Build regex
	regex := regexp.MustCompile(`stream\s+of`)

	// Find matches
	matchedData := FindReader(context.Background(), regex, 40, stream)
	matches := 0
	for range matchedData {
		matches++
		// fmt.Println(string(match))
	}
	if matches != 1 {
		t.Errorf("Did not find correct number of matches")
	}

	// Test for duplicates

	// Create string
	data = `test test test test test test test test test`
	stream = strings.NewReader(data)

	// Build regex
	regex = regexp.MustCompile(`test`)

	// Find matches
	matchedData = FindReader(context.Background(), regex, 5, stream)
	matches = 0
	for match := range matchedData {
		matches++
		if match != "test" {
			t.Errorf("Invalid match: %s", match)
		}
		// fmt.Println(string(match))
	}
	if matches != 9 {
		t.Errorf("Did not find correct number of matches")
	}
}

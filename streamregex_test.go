package streamregex

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestFindReader(t *testing.T) {
	// Create string
	data := `0123456789this is a stream    of data with lots of trailing information`
	stream := strings.NewReader(data)

	// Build regex
	regexInt := regexp.MustCompile(`stream\s+of`)
	regex := NewRegex(regexInt)
	// Use a ring buffer such that we cut off our match, but will get it in the overlap
	regex.RingBufferSize = 22
	regex.RingBufferOverlap = 5

	// Find matches
	matchedData := regex.FindReader(context.Background(), stream)
	matches := 0
	for match := range matchedData {
		matches++
		fmt.Println(string(match))
	}
	if matches != 1 {
		t.Errorf("Did not find correct number of matches")
	}
}

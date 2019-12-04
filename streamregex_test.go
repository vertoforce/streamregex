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
	data := `this is a stream    of data`
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
	if matches != 1 {
		t.Errorf("Did not find correct number of matches")
	}
}

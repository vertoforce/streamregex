package streamregex

import (
	"context"
	"regexp"
	"strings"
	"testing"
)

func TestFindReader(t *testing.T) {
	tests := []struct {
		input       string
		match       string
		regex       *regexp.Regexp
		maxMatchLen int
		numMatchs   int
	}{
		{
			`0123456789this is a stream    of data with lots of trailing information`,
			`stream    of`,
			regexp.MustCompile(`stream\s+of`),
			40,
			1,
		},
		{
			`test test test test test test test test test`,
			`test`,
			regexp.MustCompile(`test`),
			4,
			9,
		},
		{
			`cut off match`,
			`f ma`,
			regexp.MustCompile(`f ma`),
			4,
			1,
		},
	}

	for _, test := range tests {
		// Find matches
		matchedData := FindReader(context.Background(), test.regex, test.maxMatchLen, strings.NewReader(test.input))
		matches := 0
		for match := range matchedData {
			matches++
			if match != test.match {
				t.Errorf("Invalid match: %s", match)
			}
		}
		if matches != test.numMatchs {
			t.Errorf("Did not find correct number of matches")
		}
	}
}

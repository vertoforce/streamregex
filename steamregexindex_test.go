package streamregex

import (
	"context"
	"regexp"
	"strings"
	"testing"
)

func TestFindReaderIndex(t *testing.T) {
	tests := []struct {
		input       string
		match       string
		regex       *regexp.Regexp
		maxMatchLen int
		numMatchs   int
		locations   [][]int
	}{
		{
			`0123456789this is a stream    of data with lots of trailing information`,
			`stream    of`,
			regexp.MustCompile(`stream\s+of`),
			40,
			1,
			[][]int{{20, 32}},
		},
		{
			`test test test test test test test test test`,
			`test`,
			regexp.MustCompile(`test`),
			4,
			9,
			[][]int{{0, 4}, {5, 9}, {10, 14}, {15, 19}, {20, 24}, {25, 29}, {30, 34}, {35, 39}, {40, 44}},
		},
		{
			`cut off match`,
			`f ma`,
			regexp.MustCompile(`f ma`),
			4,
			1,
			[][]int{{6, 10}},
		},
	}

	for _, test := range tests {
		// Find matches
		matchedData, locationsData := FindReaderIndex(context.Background(), test.regex, test.maxMatchLen, strings.NewReader(test.input))
		matches := 0
		for match := range matchedData {
			location := <-locationsData
			matches++
			if match != test.match {
				t.Errorf("Invalid match: %s", match)
			}
			if location[0] != test.locations[matches-1][0] && location[1] != test.locations[matches-1][1] {
				t.Errorf("Invalid match location: %d", location)
			}
		}
		if matches != test.numMatchs {
			t.Errorf("Did not find correct number of matches")
		}
	}
}

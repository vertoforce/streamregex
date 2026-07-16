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

	// The caller owns the channel. Run FindReader in a goroutine and close the channel
	// once it returns so the range below terminates.
	matches := make(chan string)
	go func() {
		FindReader(context.Background(), regex, 100, stream, matches)
		close(matches)
	}()

	// Find matches
	for match := range matches {
		fmt.Println(match)
	}

	// Output: stream    of
}

func ExampleFindReaderIndex() {
	// Create string
	data := `0123456789this is a stream    of data with lots of trailing information`
	stream := strings.NewReader(data)

	// Build regex
	regex := regexp.MustCompile(`stream\s+of`)

	// The caller owns both channels. The index is sent just before its match, so the
	// index channel is buffered by 1 to avoid a deadlock. Run FindReaderIndex in a
	// goroutine and close both channels once it returns.
	matches := make(chan string)
	indexes := make(chan []int, 1)
	go func() {
		FindReaderIndex(context.Background(), regex, 100, stream, matches, indexes)
		close(matches)
		close(indexes)
	}()

	// Find matches and indexes
	for match := range matches {
		fmt.Println(match, <-indexes)
	}

	// Output: stream    of [20 32]
}

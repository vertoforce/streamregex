// Package streamregex allows you to get the matched data of a regex on a io.Reader stream
package streamregex

import (
	"bufio"
	"context"
	"io"
	"regexp"
)

// SplitRegex takes a regex and returns a split function that will find that regex in a byte slice
func SplitRegex(re *regexp.Regexp, maxMatchLength int) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, io.EOF
		}
		if loc := re.FindIndex(data); loc != nil {
			return loc[1], data[loc[0]:loc[1]], nil
		}
		if atEOF {
			return 0, nil, io.EOF
		}
		if len(data) >= maxMatchLength {
			return maxMatchLength, nil, nil
		}
		return 0, nil, nil
	}
}

// FindReader return channel of matched []byte from reader
func FindReader(ctx context.Context, r *regexp.Regexp, maxMatchLength int, reader io.Reader) chan string {
	allMatches := make(chan string)

	buf := make([]byte, maxMatchLength)

	go func() {
		defer close(allMatches)

		scanner := bufio.NewScanner(reader)
		scanner.Buffer(buf, maxMatchLength)
		scanner.Split(SplitRegex(r, maxMatchLength))
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case allMatches <- scanner.Text():
			}
		}
	}()

	return allMatches
}

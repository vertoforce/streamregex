package streamregex

import (
	"bufio"
	"context"
	"io"
	"regexp"
)

// SplitRegexIndex takes a regex and a channel for sending the locations at wich the matches are found; and returns a split function that will find that regex in a byte slice
func SplitRegexIndex(re *regexp.Regexp, maxMatchLength int, indexChannel chan []int) bufio.SplitFunc {
	var byteCount int

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, io.EOF
		}

		if loc := re.FindIndex(data); loc != nil {
			absLoc := make([]int, 2)
			absLoc[0] = loc[0] + byteCount
			absLoc[1] = loc[1] + byteCount
			indexChannel <- absLoc
			byteCount = byteCount + loc[1]
			return loc[1], data[loc[0]:loc[1]], nil
		}
		if atEOF {
			return 0, nil, io.EOF
		}
		if len(data) >= maxMatchLength {
			var advance = len(data) - maxMatchLength
			byteCount = byteCount + advance
			return advance, nil, nil
		}
		return 0, nil, nil
	}
}

// FindReaderIndex returns a channel of matched []byte from the reader and a channel of the locations at wich the matches where found
// This function will allocate maxMatchLength*2 bytes of memory
func FindReaderIndex(ctx context.Context, r *regexp.Regexp, maxMatchLength int, reader io.Reader) (chan string, chan []int) {
	allMatches := make(chan string)
	allIndexes := make(chan []int, 1) // Avoid deadlock because if location is not read it blocks SplitRegexIndex from returning the match

	buf := make([]byte, maxMatchLength*2)

	go func() {
		defer close(allMatches)

		scanner := bufio.NewScanner(reader)
		scanner.Buffer(buf, maxMatchLength)
		scanner.Split(SplitRegexIndex(r, maxMatchLength, allIndexes))
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case allMatches <- scanner.Text():
			}
		}
	}()

	return allMatches, allIndexes
}

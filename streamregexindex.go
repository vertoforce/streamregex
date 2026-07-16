package streamregex

import (
	"bufio"
	"context"
	"io"
	"regexp"
)

// SplitRegexIndex takes a regex and a channel for sending the locations at which the matches are found; and returns a split function that will find that regex in a byte slice
func SplitRegexIndex(re *regexp.Regexp, maxMatchLength int, indexChannel chan<- []int) bufio.SplitFunc {
	var byteCount int

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, io.EOF
		}

		if loc := re.FindIndex(data); loc != nil {
			absLoc := []int{loc[0] + byteCount, loc[1] + byteCount}
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

// FindReaderIndex scans reader for matches of re. It sends each match on the caller-owned
// matches channel and its location ([start, end] absolute byte offsets) on the caller-owned
// indexes channel. It runs synchronously (no internal goroutine) and returns when the reader
// is exhausted, ctx is canceled, or a read error occurs.
//
// The caller owns both channels: FindReaderIndex never closes them. Typically you launch
// FindReaderIndex in a goroutine, then close both channels once it returns and range over
// them from another goroutine (see ExampleFindReaderIndex).
//
// The caller MUST drain BOTH channels. Each index is sent immediately before its
// corresponding match, so give the index channel a buffer of 1 (or receive the index
// before the match) — otherwise the scanner blocks sending the index while the caller
// waits for the match and the whole thing deadlocks.
//
// If ctx is canceled, FindReaderIndex returns ctx.Err(). Otherwise it returns any error
// from the underlying reader (a nil return means a clean EOF). This function will allocate
// maxMatchLength*2 bytes of memory.
func FindReaderIndex(ctx context.Context, re *regexp.Regexp, maxMatchLength int, reader io.Reader, matches chan<- string, indexes chan<- []int) error {
	buf := make([]byte, maxMatchLength*2)

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(buf, maxMatchLength)
	scanner.Split(SplitRegexIndex(re, maxMatchLength, indexes))
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case matches <- scanner.Text():
		}
	}

	return scanner.Err()
}

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
			return len(data) - maxMatchLength, nil, nil
		}
		return 0, nil, nil
	}
}

// FindReader scans reader for matches of re and sends each match on the caller-owned
// matches channel. It runs synchronously (no internal goroutine) and returns when the
// reader is exhausted, ctx is canceled, or a read error occurs.
//
// The caller owns matches: FindReader never closes it. Typically you launch FindReader
// in a goroutine, then close matches once it returns and range over it from another
// goroutine (see ExampleFindReader).
//
// If ctx is canceled, FindReader returns ctx.Err(). Otherwise it returns any error from
// the underlying reader (a nil return means a clean EOF). This function will allocate
// maxMatchLength*2 bytes of memory.
func FindReader(ctx context.Context, re *regexp.Regexp, maxMatchLength int, reader io.Reader, matches chan<- string) error {
	buf := make([]byte, maxMatchLength*2)

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(buf, maxMatchLength)
	scanner.Split(SplitRegex(re, maxMatchLength))
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case matches <- scanner.Text():
		}
	}

	return scanner.Err()
}

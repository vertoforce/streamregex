// Package streamregex allows you to get the matched data of a regex on a io.Reader stream
package streamregex

import (
	"bytes"
	"context"
	"io"
	"regexp"
)

// Defaults for ring buffer
const (
	DefaultRingBufferSize    = 1024 * 1024 // 1MB
	DefaultRingBufferOverlap = 1024        // 1KB
)

// Regex just a wrapper around a regex
type Regex struct {
	Regex             *regexp.Regexp
	RingBufferSize    int
	RingBufferOverlap int
}

// NewRegex Create regex from built in regex package
func NewRegex(regex *regexp.Regexp) *Regex {
	r := &Regex{
		Regex:             regex,
		RingBufferSize:    DefaultRingBufferSize,
		RingBufferOverlap: DefaultRingBufferOverlap,
	}

	return r
}

// FindReaderString return channel of matched strings from reader
func (r *Regex) FindReaderString(ctx context.Context, reader io.Reader) chan string {
	ret := make(chan string)

	go func() {
		defer close(ret)

		for match := range r.FindReader(ctx, reader) {
			select {
			case ret <- string(match):
			case <-ctx.Done():
				return
			}
		}
	}()

	return ret
}

// FindReader return channel of matched []byte from reader
func (r *Regex) FindReader(ctx context.Context, reader io.Reader) chan []byte {
	allMatches := make(chan []byte)

	go func() {
		defer close(allMatches)

		// Read and scan in chunks with some overlap
		var nextOverlap []byte
		for {
			// Read into ring buffer
			buf := &bytes.Buffer{}

			// Write last overlap
			buf.Write(nextOverlap)

			// Read from stream until buffer is full len==(RingBufferSize)
			n, err := io.CopyN(buf, reader, int64(r.RingBufferSize-len(nextOverlap)))
			if err != nil && err != io.EOF {
				// Real error reading, just break off
				break
			}

			// Scan what we actually read
			matches := r.Regex.FindAll(buf.Bytes()[0:n], -1)
			for _, match := range matches {
				select {
				case allMatches <- match:
				case <-ctx.Done():
					return
				}
			}

			// See if we are done and should stop
			if int(n) < r.RingBufferSize || err == io.EOF {
				// We are done
				break
			} else {
				// We should keep going, copy the end of this buffer in to start of next buffer
				nextOverlap = buf.Bytes()[buf.Len()-1-r.RingBufferOverlap : buf.Len()]
			}
		}
	}()

	return allMatches
}

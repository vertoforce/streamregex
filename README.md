# Streamregex

[![Go Report Card](https://goreportcard.com/badge/github.com/vertoforce/streamregex)](https://goreportcard.com/report/github.com/vertoforce/streamregex)
[![Documentation](https://godoc.org/github.com/vertoforce/streamregex?status.svg)](https://godoc.org/github.com/vertoforce/streamregex)

Go does not let you find the _matched data_ of a regex from a stream of data.  They let you find [the first position of a match](https://godoc.org/regexp#Regexp.FindReaderIndex), but not the data itself.

Streamregex sends you the _matched data_ of a regex on a io.Reader stream over a channel.

## Usage

`FindReader` sends each match on a channel you own. It runs synchronously and returns an
error (`nil` on clean EOF, otherwise a read or context error), so the caller owns error
handling, concurrency, and the channel's lifecycle. Run it in a goroutine and close the
channel once it returns:

```go
// Create string
data := `0123456789this is a stream    of data with lots of trailing information`
stream := strings.NewReader(data)

// Build regex
regex := regexp.MustCompile(`stream\s+of`)

// The caller owns the channel: run FindReader in a goroutine and close it when done.
matches := make(chan string)
go func() {
    FindReader(context.Background(), regex, 100, stream, matches)
    close(matches)
}()

for match := range matches {
    fmt.Println(match)
}

// Output: stream    of
```

### With match locations

`FindReaderIndex` additionally reports each match's `[start, end]` byte offsets on a second
channel. The index is sent just before its match, so the caller MUST drain both channels;
buffer the index channel by 1 (or receive the index before the match) to avoid a deadlock:

```go
matches := make(chan string)
indexes := make(chan []int, 1) // buffered: index is sent just before its match
go func() {
    FindReaderIndex(context.Background(), regex, 100, stream, matches, indexes)
    close(matches)
    close(indexes)
}()

for match := range matches {
    fmt.Println(match, <-indexes)
}

// Output: stream    of [20 32]
```

## How it works

We use a custom [SplitFunc](https://golang.org/pkg/bufio/#SplitFunc) to split the reader into each regex match.  Normally for a SplitFunc it will keep reading more and more data into the buffer until it finds a match. To avoid pulling all the reader data into memory, the function accepts a `maxMatchLength` if you know the maximum match length of a match.

Note that we need to allocate a `maxMatchLength*2` bytes of memory to successfully scan the reader for matches.

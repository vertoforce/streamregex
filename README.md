# Streamregex

[![Go Report Card](https://goreportcard.com/badge/github.com/vertoforce/streamregex)](https://goreportcard.com/report/github.com/vertoforce/streamregex)
[![Documentation](https://godoc.org/github.com/vertoforce/streamregex?status.svg)](https://godoc.org/github.com/vertoforce/streamregex)

Go does not let you find the _matched data_ of a regex from a stream of data.  They let you find [the first position of a match](https://godoc.org/regexp#Regexp.FindReaderIndex), but not the data itself.

Streamregex allows you to get a channel of the _matched data_ of a regex on a io.Reader stream.

## Usage

```go
// Create string
data := `0123456789this is a stream    of data with lots of trailing information`
stream := strings.NewReader(data)

// Build regex
regex := regexp.MustCompile(`stream\s+of`)

// Find matches
matchedData := FindReader(context.Background(), regex, 100, stream)
for match := range matchedData {
    fmt.Println(string(match))
}

// Output: stream    of
```

## How it works

We use a custom [SplitFunc](https://golang.org/pkg/bufio/#SplitFunc) to split the reader into each regex match.  Normally for a SplitFunc it will keep reading more and more data into the buffer until it finds a match. To avoid pulling all the reader data into memory, the function accepts a `maxMatchLength` if you know the maximum match length of a match.

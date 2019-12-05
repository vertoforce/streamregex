# streamregex

streamregex allows you to get the matched data of a regex on a io.Reader stream.

## Usage

```go
// Build stream
streamData = `this is a stream of data`
stream := strings.NewReader(streamData)

// Build regex
regexInt := regexp.MustCompile(`stream of`)
regex := Regex(regexInt)
regex.RingBufferSize = 1024*1024 // 1MB (default)
regex.RingBufferOverlap = 1024 // 1KB (default)

matchedData := regex.FindReader(stream)
for match := range matchedData {
    fmt.Println(match)
}
```

Outputs: `stream of`

## How it works

This basically uses a sliding window buffer to scan parts of the input stream.  You can configure this size with
`regex.RingBufferSize` and `regex.RingBufferOverlap` based on the expected length of your rules.

Note that to avoid duplicate rule matches, it doesn't add matches that are exactly the same as the last match.  So if you are
expecting multiple matches in a row, you may not see them.

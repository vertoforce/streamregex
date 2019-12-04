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

matchedData := regex.FindReader(stream)
for match := range matchedData {
    fmt.Println(match)
}
```

Outputs: `stream of`

# qst

`qst` is an *http.Request builder.

Install with:
```shell script
$ go get github.com/broothie/qst
```

`qst` uses an options pattern:

```go
request, err := qst.NewPatch("http://example.com",   // New PATCH request
    qst.Bearer("some-token-here"),                   // Authorization header
    qst.QueryValue("key", "value"),                  // Query param
    qst.BodyJSON(map[string]string{"key": "value"}), // JSON body
)
```

It can also be used to fire requests:

```go
response, err := qst.Patch("http://example.com",     // Send PATCH request
    qst.Bearer("some-token-here"),                   // Authorization header
    qst.QueryValue("key", "value"),                  // Query param
    qst.BodyJSON(map[string]string{"key": "value"}), // JSON body
)
```

The options pattern allows for easily defining commonly used options:

```go
func query(before time.Time) qst.Option {
    return OptionList{
        qst.Authorization("some-token"),
        qst.QueryValue("created_at", fmt.Sprintf(">=%s", before.Format(time.RFC3339))),
    }
}

func main() {
    qst.Get("http://example.com", query(time.Date(1993, time.March, 19, 0, 0, 0, 0, time.UTC))
}
```

If you wish to use an existing *http.Client:

```go
func main() {
    client := &http.Client{Timeout: 3 * time.Second}
    response, err := qst.WithClient(client).Post("http://example.com", qst.Bearer("token"), qst.BodyJSON(payload))
}
```

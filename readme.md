# qst

[![Go Reference](https://pkg.go.dev/badge/github.com/broothie/qst.svg)](https://pkg.go.dev/github.com/broothie/qst)
[![Go Report Card](https://goreportcard.com/badge/github.com/broothie/qst)](https://goreportcard.com/report/github.com/broothie/qst)

`qst` is an `*http.Request` builder. "qst" is short for "quest", which is part of the word "request".

## Installation

```shell script
$ go get github.com/broothie/qst
```

## Documentation

Detailed documentation can be found at [pkg.go.dev](https://pkg.go.dev/github.com/broothie/qst).

## Usage

`qst` uses an options pattern to build `*http.Request` objects:
```go
request, err := qst.NewPatch("http://example.com",   // New PATCH request
    qst.Bearer("some-token-here"),                   // Authorization header
    qst.QueryValue("key", "value"),                  // Query param
    qst.BodyJSON(map[string]string{"key": "value"}), // JSON body
)
```

Documentation for all available options can be found [here](https://pkg.go.dev/github.com/broothie/qst#Option).

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
func createdAfter(after time.Time) qst.Option {
    return qst.QueryValue("created_at", fmt.Sprintf(">=%s", after.Format(time.RFC3339)))
}

func main() {
    response, err := qst.Get("http://example.com", createdAfter(time.Now().Add(-24 * time.Hour)))
}
```

If you wish to use an existing `*http.Client`:
```go
func makeCall(token string, payload map[string]interface{}) {
    client := &http.Client{Timeout: 3 * time.Second}
    response, err := qst.WithClient(client).Post("http://example.com", qst.Bearer(token), qst.BodyJSON(payload))
}
```

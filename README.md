# qst

[![Go Reference](https://pkg.go.dev/badge/github.com/broothie/qst.svg)](https://pkg.go.dev/github.com/broothie/qst)
[![Go Report Card](https://goreportcard.com/badge/github.com/broothie/qst)](https://goreportcard.com/report/github.com/broothie/qst)
[![codecov](https://codecov.io/gh/broothie/qst/branch/main/graph/badge.svg?token=CVMUN8Y9FV)](https://codecov.io/gh/broothie/qst)
[![gosec](https://github.com/broothie/qst/actions/workflows/gosec.yml/badge.svg)](https://github.com/broothie/qst/actions/workflows/gosec.yml)

`qst` is an `*http.Request` builder. "qst" is short for "quest", which is part of the word "request".

## Installation

```shell script
go get github.com/broothie/qst
```

## Documentation

Detailed documentation can be found at [pkg.go.dev](https://pkg.go.dev/github.com/broothie/qst).

A list of all available options can be found [here](https://pkg.go.dev/github.com/broothie/qst#Option).

## Usage

`qst` uses an options pattern to build `*http.Request` objects:

```go
request, err := qst.NewPatch("https://breakfast.com/api", // New PATCH request
    qst.WithBearerAuth("c0rNfl@k3s"),                         // Authorization header
    qst.WithPath("/cereals", cerealID),                       // Query param
    qst.WithBodyJSON(map[string]string{"name": "Life"}),      // JSON body
)
```

It can also be used to fire requests:

```go
response, err := qst.Patch("https://breakfast.com/api", // Send PATCH request
    qst.WithBearerAuth("c0rNfl@k3s"),                       // Authorization header
    qst.WithPath("/cereals", cerealID),                     // Query param
    qst.WithBodyJSON(map[string]string{"name": "Life"}),    // JSON body
)
```

The options pattern makes it easy to define custom options:

```go
func createdSinceYesterday() option.Option[*http.Request] {
    yesterday := time.Now().Add(-24 * time.Hour)
    return qst.WithQuery("created_at", fmt.Sprintf(">%s", yesterday.Format(time.RFC3339)))
}

func main() {
    response, err := qst.Get("https://breakfast.com/api",
        qst.WithBearerAuth("c0rNfl@k3s"),
        qst.WithPath("/cereals"),
        createdSinceYesterday(),
    )
}
```

You can also combine multiple options using `option.NewOptions()` to create reusable defaults:

```go
func myDefaults() option.Option[*http.Request] {
    return option.NewOptions(
        qst.WithURL("https://breakfast.com/api"),
        qst.WithBearerAuth("c0rNfl@k3s"),
    )
}
```

## All Available Options

```go
request, err := qst.New(
    // Use a *url.URL
    qst.WithRawURL(parsedURL),

    // Use a URL string
    qst.WithURL("https://api.example.com"),

    // Set scheme
    qst.WithScheme("https"),

    // Set host
    qst.WithHost("api.example.com"),

    // Build path from segments
    qst.WithPath("/users", userID, "/posts"),

    // Use *url.Userinfo
    qst.WithUser(userInfo),

    // Username only
    qst.WithUsername("admin"),

    // Username and password
    qst.WithUserPassword("admin", "secret"),

    // Single query param
    qst.WithQuery("page", "1"),

    // Multiple query params
    qst.WithQueries(qst.Queries{
        "page":  {"1"},
        "limit": {"10"},
    }),

    // Single header
    qst.WithHeader("X-API-Key", "secret"),

    // Multiple headers
    qst.WithHeaders(qst.Headers{
        "X-API-Key":    {"secret"},
        "X-Client-ID":  {"app123"},
    }),

    // Accept header
    qst.WithAcceptHeader("application/json"),

    // Content-Type header
    qst.WithContentTypeHeader("application/xml"),

    // Referer header
    qst.WithRefererHeader("https://example.com"),

    // User-Agent header
    qst.WithUserAgentHeader("MyApp/1.0"),

    // Authorization header
    qst.WithAuthorizationHeader("Custom token"),

    // Basic auth
    qst.WithBasicAuth("user", "pass"),

    // Token auth
    qst.WithTokenAuth("abc123"),

    // Bearer token
    qst.WithBearerAuth("jwt_token_here"),

    // Add cookie
    qst.WithCookie(&http.Cookie{
        Name:  "session",
        Value: "abc123",
    }),

    // Set context
    qst.WithContext(ctx),

    // Add context value
    qst.WithContextValue("userID", 123),

    // io.ReadCloser
    qst.WithBody(readCloser),

    // io.Reader
    qst.WithBodyReader(reader),

    // Byte slice
    qst.WithBodyBytes([]byte("data")),

    // String
    qst.WithBodyString("plain text"),

    // URL-encoded form
    qst.WithBodyForm(qst.Form{
        "username": {"john"},
        "email":    {"john@example.com"},
    }),

    // JSON body
    qst.WithBodyJSON(map[string]interface{}{
        "name": "John",
        "age":  30,
    }),

    // XML body
    qst.WithBodyXML(struct{
        Name string `xml:"name"`
        Age  int    `xml:"age"`
    }{
        Name: "John",
        Age:  30,
    }),

    // Dump request to writer
    qst.WithDump(os.Stdout),
)
```

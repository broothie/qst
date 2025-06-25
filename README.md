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

### Creating Default Options

If you need to reuse a set of default options across multiple requests, you can create a function that returns multiple options:

```go
func myDefaults() option.Option[*http.Request] {
    return option.NewOptions(
        qst.WithURL("https://breakfast.com/api"),
        qst.WithBearerAuth("c0rNfl@k3s"),
    )
}

func main() {
    response, err := qst.Patch("https://breakfast.com/api",
        myDefaults(),
        qst.WithPath("/cereals", cerealID),
        qst.WithBodyJSON(map[string]interface{}{
            "name": "Golden Grahams",
        }),
    )
}
```

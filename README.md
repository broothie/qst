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
    qst.BearerAuth("c0rNfl@k3s"),                         // Authorization header
    qst.Path("/cereals", cerealID),                       // Query param
    qst.BodyJSON(map[string]string{"name": "Life"}),      // JSON body
)
```

It can also be used to fire requests:

```go
request, err := qst.Patch("https://breakfast.com/api", // Send PATCH request
    qst.BearerAuth("c0rNfl@k3s"),                      // Authorization header
    qst.Path("/cereals", cerealID),                    // Query param
    qst.BodyJSON(map[string]string{"name": "Life"}),   // JSON body
)
```

The options pattern makes it easy to define custom options:

```go
func createdSinceYesterday() qst.Option {
    yesterday := time.Now().Add(-24 * time.Hour)
    return qst.Query("created_at", fmt.Sprintf(">%s", yesterday.Format(time.RFC3339)))
}

func main() {
    response, err := qst.Get("https://breakfast.com/api",
        qst.BearerAuth("c0rNfl@k3s"),
        qst.Path("/cereals"),
        createdSinceYesterday(),
    )
}
```

### qst.Client

This package also includes a `Client`, which can be outfitted with a set of default options:

```go
client := qst.NewClient(http.DefaultClient,
    qst.URL("https://breakfast.com/api"),
    qst.BearerAuth("c0rNfl@k3s"), 
)

response, err := client.Patch(
    // qst.URL("https://breakfast.com/api"), // Not necessary, included via client
    // qst.BearerAuth("c0rNfl@k3s"),         // Not necessary, included via client
    qst.Path("/cereals", cerealID),
    qst.BodyJSON(map[string]interface{}{
        "name": "Golden Grahams",
    }),
)
```

### qst.OptionFunc

`OptionFunc` can be used to add a custom function which is run during request creation:

```go
client := qst.NewClient(http.DefaultClient,
    qst.OptionFunc(func(request *http.Request) (*http.Request, error) {
      token := dynamicallyGetBearerTokenSomehow()
      return qst.BearerAuth(token).Apply(request)
    }),
)
```

// Package qst is an *http.Request builder. It uses a variadic options pattern:
//     request, err := qst.NewPatch("http://example.com",   // New PATCH request
//         qst.BearerAuth("some-token-here"),               // Authorization header
//         qst.QueryValue("key", "value"),                  // Query param
//         qst.BodyJSON(map[string]string{"key": "value"}), // JSON body
//     )
//
// It can also be used to fire requests:
//     response, err := qst.Patch("http://example.com",     // Send PATCH request
//         qst.BearerAuth("some-token-here"),               // Authorization header
//         qst.QueryValue("key", "value"),                  // Query param
//         qst.BodyJSON(map[string]string{"key": "value"}), // JSON body
//     )
//
// The variadic option pattern allows for easily defining commonly used options:
//     func query(before time.Time) qst.Option {
//         return Pipeline{
//             qst.Authorization("some-token"),
//             qst.QueryValue("created_at", fmt.Sprintf(">=%s", before.Format(time.RFC3339))),
//         }
//     }
//
//     func main() {
//         qst.Get("http://example.com", query(time.Date(1993, time.March, 19, 0, 0, 0, 0, time.UTC))
//     }
//
// If you wish to use an existing *http.Client:
//     func main() {
//         client := &http.Client{Timeout: 3 * time.Second}
//         response, err := qst.WithClient(client).Post("http://example.com", qst.BearerAuth("token"), qst.BodyJSON(payload))
//     }
//
package qst

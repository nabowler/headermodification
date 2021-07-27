# headermodification

A simple library to control request headers when you control the http.Client, but not necessarily the request.

## Use

```go
package main

import (
    "net/http"

    "github.com/nabowler/headermodification"
)

func main() {
    transport := headermodification.Transport{
        // the headers have been pre-canonicalized for this example
        Set: http.Header{"User-Agent": []string{"my custom useragent"}},
        Add: http.Header{"X-Custom-1": []string{"some custom header"}},
    }

    client := http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }

    api(client)
}


func api(client http.Client) {
    req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "https://example.com", nil)
    if err != nil {
        // TODO
    }

    resp, err := client.Do(req)
    if err != nil {
        // TODO
    }
    defer resp.Body.Close()
}
```
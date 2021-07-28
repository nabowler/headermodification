# headermodification

A simple library to control request headers when you control the http.Client, but not necessarily the request.

This library provides an `http.RoundTripper` implementation that will Set and Add headers on the outbound request.

## Use

The following snippet configures the `Transport` on an `http.Client` and then provides that client to an API layer that performs the requests. Because the API layer controls the requests, we would normally not be able to control the outbound headers.

```go
package main

import (
    "net/http"

    "github.com/nabowler/headermodification"
)

func main() {
    transport := headermodification.Transport {
        Set: http.Header{"User-Agent": []string{"my custom useragent"}},
        Add: http.Header{"X-Custom-1": []string{"some custom header"}},
    }

    client := http.Client {
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

### Global

The following snippet configures the `Transport` as `http.DefaultTransport`. This is necessary because the API controls both the `http.Client` and the requests, so normally
we would have no control over the requests. By hijacking the `http.DefaultTransport` global variable, most HTTP requests from the application will have their header's modified.

This is not recommended unless necessary and should only ever be configured as such within `func main()`. This pattern must be used with caution since this may affect all HTTP requests from your application and not just those of a single API.

```go
package main

import (
    "net/http"

    "github.com/nabowler/headermodification"
)

func main() {
    transport := headermodification.Transport {
        Base: http.DefaultTransport,
        Set: http.Header{"User-Agent": []string{"my custom useragent"}},
        Add: http.Header{"X-Custom-1": []string{"some custom header"}},
    }

    http.DefaultTransport = transport

    api()
}


func api() {
    client := http.Client {
        Timeout:   30 * time.Second,
    }

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
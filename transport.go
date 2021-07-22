package headermodification

import (
	"net/http"
)

// Transport will modify the headers of the out-bound request using the
// values in Set and Add before using Base to perform the RoundTrip.
// Values in Set will be applied before values in Add.
type Transport struct {
	// Base is the RoundTripper that will be used after the request headers
	// are modified. If nil, http.DefaultTransport will be used.
	Base http.RoundTripper
	// Set contains the headers that will be overwritten using `http.Header.Set`.
	Set http.Header
	// Add contains the headers that will be appended using `http.Header.Add`.
	Add http.Header
}

// RoundTrip modifies the headers of the out-bound request using the
// values in Set and Add before using Base to perform the RoundTrip.
// If Base is not set, http.DefaultTransport will be used.
// If both Set and Add are empty, the request is sent unmodified.
// If Set or Add have values, the request is cloned before modification,
// and the clone will be modified (per the RoundTripper contract) and sent.
func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	set := t.Set
	add := t.Add

	// if there's nothing configured, we can just submit the request as is
	if len(set) == 0 && len(add) == 0 {
		return base.RoundTrip(req)
	}

	req2 := cloneRequest(req) // per RoundTripper contract

	ModifyHeaders(req2, set, add)

	return base.RoundTrip(req2)
}

// ModifyHeaders directly modifies the supplied request's headers.
// Existing headers will be overwritten with the values in set.
// Existing headers, including those set here, will be appeneded to
// with the values in add.
// Headers that are not in set or add will remain unmodified.
func ModifyHeaders(req *http.Request, set, add http.Header) {
	for k, vs := range set {
		if len(vs) == 0 {
			continue
		}
		// set the first value to guarantee the header is cleared
		req.Header.Set(k, vs[0])
		if len(vs) > 1 {
			// add any values after the first
			for _, v := range vs[1:] {
				req.Header.Add(k, v)
			}
		}
	}

	for k, vs := range add {
		if len(vs) == 0 {
			continue
		}
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
}

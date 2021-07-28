package headermodification

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestClone(t *testing.T) {
	type ctxkey string
	const key ctxkey = "key"
	const value = "value"

	ctx := context.WithValue(context.Background(), key, value)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", io.NopCloser(strings.NewReader("body")))
	if err != nil {
		t.Fatal(err.Error())
	}

	clone := cloneRequest(req)

	if clone == nil {
		t.Fatal("clone is nil")
	}

	if req == clone {
		// clone must not point to the same object
		t.Error("req and clone point to the same request")
	}

	if req.Method != clone.Method {
		t.Error("methods are different")
	}
	if req.URL != clone.URL {
		t.Error("urls are different")
	}
	if req.Body != clone.Body {
		t.Error("bodies are different")
	}
	if !reflect.DeepEqual(req.Header, clone.Header) {
		t.Error("headers are not deeply equal")
	}
	if req.Context() != clone.Context() {
		t.Error("contexts are different")
	}

	// many of the above checks are probably also covered by this
	// Note: this check passes with a nil body, or a io.NopCloser body,
	// but fails with a non-nil strings Reader or bytes Reader/Buffer because
	// GetBody() gets set for these body types, and functions cannot be compared
	// between objects
	if !reflect.DeepEqual(req, clone) {
		t.Errorf("req and clone are not deeply equal")
	}
}

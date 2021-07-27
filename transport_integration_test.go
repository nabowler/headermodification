package headermodification_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/nabowler/headermodification"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping because the -test.short flag is set")
	}

	const useragent = "headermodification_testintegration"
	const custom = "this is a custom value"

	transport := headermodification.Transport{
		// the headers have been pre-canonicalized for this test
		Add: http.Header{"X-Custom-1": []string{custom}},
		Set: http.Header{"User-Agent": []string{useragent}},
	}

	client := http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://httpbin.org/headers", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("User-Agent", "this should be overwritten")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer resp.Body.Close()

	var response struct {
		Headers map[string]string `json:"headers"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatal(err.Error())
	}

	ua, ok := response.Headers["User-Agent"]
	if ok {
		if useragent != ua {
			t.Fatalf("Unexpected user-agent header: %s", ua)
		}
	} else {
		t.Fatal("User-Agent header not found in response")
	}

	c, ok := response.Headers["X-Custom-1"]
	if ok {
		if custom != c {
			t.Fatalf("Unexpected X-Custom-1 header: %s", c)
		}
	} else {
		t.Fatal("X-Custom-1 header not found in response")
	}

	if t.Failed() {
		t.Logf("Response: %+v", response)
	}
}

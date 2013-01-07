package subscribe

import (
	"flag"
	"testing"
)

var testClientID = flag.String("clientID", "", "Enable tests with this client ID")

func TestSubscribe(t *testing.T) {
	if *testClientID == "" {
		t.Logf("Skipping test (no --clientID specified)")
		return
	}

	urls := make(chan string)
	go func() {
		if err := Subscribe(*testClientID, urls); err != nil {
			t.Errorf("subscribe: %s", err)
		}
	}()

	for url := range urls {
		t.Logf("got %q", url)
	}
}

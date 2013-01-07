// Package subscribe allows for subscriptions to the gopaste service
package subscribe

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/kylelemons/gaechannel"
)

var PasteServer = url.URL{
	Scheme: "http",
	Host:   "gopaste.kevlar-go-test.appspot.com",
}

// Subscribe creates a channel with the given Client ID and sends all new URLs
// on the urls channel.  Subscribe blocks until the connection fails.  If there
// is an error in setup, decoding, or when the connection fails, Subscribe
// returns the offending error.
func Subscribe(clientID string, urls chan<- string) error {
	u := PasteServer
	u.Path = "/listen/"
	u.RawQuery = url.Values{
		"clientID": {clientID},
	}.Encode()

	// Get a token
	resp, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("get token: %s", err)
	}
	rawtok, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("read token: %s", err)
	}
	token := strings.TrimSpace(string(rawtok))

	channel := gaechannel.New(PasteServer.Host, clientID, token)
	defer channel.Close()
	return channel.Stream(urls)
}

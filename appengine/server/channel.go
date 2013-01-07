package server

import (
	"io"
	"net/http"

	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"appengine/delay"
)

func init() {
	http.HandleFunc("/_ah/channel/connected/", connected)
	http.HandleFunc("/_ah/channel/disconnected/", disconnected)
	http.HandleFunc("/listen/", channelListen)
}

type client struct{}

func connected(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	clientID := r.FormValue("from")

	key := datastore.NewKey(ctx, "client", clientID, 0, nil)
	if _, err := datastore.Put(ctx, key, &client{}); err != nil {
		ctx.Errorf("put(%q): %s", clientID, err)
		return
	}
	ctx.Infof("Client %q connected", clientID)
}

func disconnected(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	clientID := r.FormValue("from")

	key := datastore.NewKey(ctx, "client", clientID, 0, nil)
	if err := datastore.Delete(ctx, key); err != nil {
		ctx.Errorf("delete(%q): %s", clientID, err)
		return
	}
	ctx.Infof("Client %q disconnected", clientID)
}

func channelListen(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	clientID := r.FormValue("clientID")
	if clientID == "" {
		http.Error(w, "no ClientID specified", http.StatusBadRequest)
		return
	}

	token, err := channel.Create(ctx, clientID)
	if err != nil {
		ctx.Infof("create(%q): %s", clientID, err)
		http.Error(w, "unable to create channel", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, token)
}

var broadcastPaste = delay.Func("broadcastPaste", func(ctx appengine.Context, url string) {
	q := datastore.NewQuery("client").KeysOnly()
	keys, err := q.GetAll(ctx, nil)
	if err != nil {
		ctx.Errorf("getall: %s", err)
		return
	}
	ctx.Infof("Found %d clients", len(keys))
	for _, client := range keys {
		clientID := client.StringID()
		if err := channel.Send(ctx, clientID, url); err != nil {
			ctx.Errorf("send(%q): %s", clientID, err)
			continue
		}
		ctx.Infof("Sent URL to %q", clientID)
	}
})

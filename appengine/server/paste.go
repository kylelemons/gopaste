package server

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
	"unicode/utf8"

	"appengine"
	"appengine/datastore"

	//"github.com/kylelemons/go-rpcgen/webrpc"
	"proto" //"github.com/kylelemons/gopaste/proto"
)

const (
	expiry   = 1 * time.Hour // Time to keep pastes
	maxsize  = 1000 * 1000   // Max size of files
	maxname  = 32            // Max length of name
	maxcount = 100           // Max pastes to keep
)

func b64sha1(data []byte) string {
	sha1 := sha1.New()
	sha1.Write(data)
	hash := sha1.Sum(nil)

	b := new(bytes.Buffer)
	b64 := base64.NewEncoder(base64.URLEncoding, b)
	b64.Write(hash)

	return b.String()
}

type server struct{}

var Server = &server{}

type paste struct {
	Origin string
	Pasted time.Time
	Data   []byte
}

func (s *server) Paste(r *http.Request, in *proto.ToPaste, out *proto.Posted) error {
	ctx := appengine.NewContext(r)

	if len(in.Data) == 0 {
		return errors.New("cowardly refusing to create zero-length paste")
	}

	if len(in.Data) > maxsize {
		return errors.New("maximum paste size exceeded")
	}

	if !utf8.Valid(in.Data) {
		return errors.New("invalid UTF-8 data (you can only paste text)")
	}

	// Sanitize name
	name := in.GetName()
	if name == "" {
		name = b64sha1(in.Data)[:10]
	} else if len(name) > maxname {
		name = name[:maxname]
	}

	ctx.Infof("Accepting %q (%d bytes)", name, len(in.Data))

	key := datastore.NewKey(ctx, "paste", name, 0, nil)
	p := &paste{
		Origin: r.RemoteAddr,
		Pasted: time.Now().UTC(),
		Data:   in.Data,
	}
	if _, err := datastore.Put(ctx, key, p); err != nil {
		ctx.Errorf("put(%q): %s", name, err)
		return err
	}

	if r.Host == "" {
		r.Host = "gopaste.kevlar-go-test.appspot.com"
	}

	// Compute the pasted URL
	outURL := (&url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   path.Join("/", name),
	}).String()
	out.Url = &outURL
	broadcastPaste.Call(ctx, outURL)
	return nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	name := r.URL.Path[1:]

	key := datastore.NewKey(ctx, "paste", name, 0, nil)

	var p paste
	if err := datastore.Get(ctx, key, &p); err != nil {
		ctx.Errorf("get(%q): %s", name, err)
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	ctx.Infof("Serving %q (%d bytes)", name, len(p.Data))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(p.Data)))
	w.Write(p.Data)
	return
}

func init() {
	proto.RegisterGoPasteWeb(Server, nil)
	http.Handle("/", Server)
}

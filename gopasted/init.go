package main

import (
	"strconv"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"flag"
	"github.com/kylelemons/gopaste/proto"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"
	_ "github.com/kylelemons/go-rpcgen/webrpc"
)

var (
	addr = flag.String("http", ":4114", "Address on which to bind for HTTP")
	base = flag.String("url", "http://paste.kylelemons.net:4114/", "Base URL on which we serve")
	expiry = flag.Duration("expiry", 1*time.Hour, "Time to keep pastes")
	maxsize = flag.Int("maxbytes", 1024*1024, "Max size of files")
	maxname = flag.Int("maxname", 32, "Max length of name")
	maxcount = flag.Int("maxcount", 100, "Max pastes to keep")
)

// URL holds the base URL on which we serve
var URL *url.URL

func b64sha1(data []byte) string {
	sha1 := sha1.New()
	sha1.Write(data)
	hash := sha1.Sum(nil)

	b := new(bytes.Buffer)
	b64 := base64.NewEncoder(base64.URLEncoding, b)
	b64.Write(hash)

	return b.String()
}

type server struct {
	lock   sync.RWMutex
	cond   sync.Cond
	pastes map[string]string
	opened map[string]time.Time
	recent []string
}

func newServer() *server {
	s := &server{
		pastes: make(map[string]string),
		opened: make(map[string]time.Time),
	}
	s.cond.L = s.lock.RLocker()
	return s
}

func (s *server) purge(delay time.Duration) {
	time.Sleep(delay)

	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.recent) == 0 {
		return
	}

	delcount := len(s.recent) - *maxcount
	oldest := time.Now().Add(-*expiry)

	i, name := 0, ""
	for i, name = range s.recent {
		if _, ok := s.pastes[name]; !ok {
			continue
		}
		if s.opened[name].Before(oldest) || i < delcount {
			log.Printf("Expiring %q (%d bytes)", name, len(s.pastes[name]))
			delete(s.pastes, name)
			delete(s.opened, name)
			continue
		}
		break
	}
	s.recent = s.recent[i:]
}

func (s *server) Paste(r *http.Request, in *proto.ToPaste, out *proto.Posted) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(in.Data) == 0 {
		return errors.New("cowardly refusing to create zero-length paste")
	}

	if len(in.Data) > *maxsize {
		return errors.New("maximum paste size exceeded")
	}

	if len(s.opened) >= *maxcount {
		log.Printf("Too many pastes; attempting purge...")
		s.lock.Unlock()
		s.purge(0)
		s.lock.Lock()
	}

	var name string
	if in.Name == nil || len(*in.Name) == 0 {
		name = b64sha1(in.Data)[:10]
	} else {
		name = *in.Name
	}
	if len(name) > *maxname {
		name = name[:*maxname]
	}

	outURL := *URL
	outURL.Path = path.Join(outURL.Path, name)
	strURL := outURL.String()

	log.Printf("Accepting %q (%d bytes)", name, len(in.Data))
	s.pastes[name] = string(in.Data)
	s.opened[name] = time.Now()
	s.recent = append(s.recent, name)

	go s.purge(*expiry)
	go s.cond.Broadcast()

	out.Url = &strURL
	return nil
}

func (s *server) Next(r *http.Request, in *proto.Empty, out *proto.Posted) error {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	log.Printf("Subscribing %q for next update...", r.RemoteAddr)

	// Always wait at least once
	s.cond.Wait()

	// Keep waiting if there are no recent ones to get
	for len(s.recent) == 0 {
		s.cond.Wait()
	}

	last := s.recent[len(s.recent)-1]
	outURL := *URL
	outURL.Path = path.Join(outURL.Path, last)

	log.Printf("Sending %q to %q", last, r.RemoteAddr)

	strURL := outURL.String()
	out.Url = &strURL
	return nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if len(r.URL.Path) == 0 {
		log.Printf("Zero-length PATH?")
		http.Error(w, "bad path", http.StatusInternalServerError)
		return
	}

	name := r.URL.Path[1:]
	paste, ok := s.pastes[name]
	if !ok {
		http.NotFound(w, r)
		return
	}

	log.Printf("Serving %q (%d bytes)", name, len(paste))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(paste)))
	io.WriteString(w, paste)
	return
}

func main() {
	flag.Parse()

	url, err := url.Parse(*base)
	if err != nil {
		log.Fatalf("bad url: %s", err)
	}
	URL = url

	s := newServer()
	proto.RegisterGoPasteWeb(s, nil)
	http.Handle("/", s)

	log.Printf("Listening on %q", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("listen/serve: %s", err)
	}
}

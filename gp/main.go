package main

import (
	"flag"
	"github.com/kylelemons/go-rpcgen/webrpc"
	"github.com/kylelemons/gopaste/proto"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

var (
	baseURL = flag.String("url", "http://gopaste.kevlar-go-test.appspot.com/", "The base URL of the GoPaste server")
	name    = flag.String("name", "", "The name of the paste (use filename or MD5 sum if not provided)")
	fname   = flag.String("f", "", "The name of a file to read (standard input if not provided)")
)

func main() {
	flag.Parse()

	url, err := url.Parse(*baseURL)
	if err != nil {
		log.Fatalf("parse url: %s", err)
	}

	paste := proto.NewGoPasteWebClient(webrpc.JSON, url)

	in, out := proto.ToPaste{}, proto.Posted{}
	file := os.Stdin
	if *fname != "" {
		if file, err = os.Open(*fname); err != nil {
			log.Fatalf("open(%q): %s", *fname, err)
		}
		*fname = filepath.Base(*fname)
		in.Name = fname
	}
	if in.Data, err = ioutil.ReadAll(file); err != nil {
		log.Fatalf("read: %s", err)
	}

	if *name != "" {
		in.Name = name
	}

	if err := paste.Paste(&in, &out); err != nil {
		log.Fatalf("paste: %s", err)
	}

	if out.Url != nil {
		log.Printf("pasted to %s", *out.Url)
	} else {
		log.Fatalf("unknown error in paste")
	}
}

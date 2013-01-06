// Code generated by protoc-gen-go.
// source: paste.proto
// DO NOT EDIT!

package proto

import proto1 "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

import "net"
import "net/rpc"
import "github.com/kylelemons/go-rpcgen/codec"
import "net/url"
import "net/http"
import "github.com/kylelemons/go-rpcgen/webrpc"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto1.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type ToPaste struct {
	Name             *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Data             []byte  `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *ToPaste) Reset()         { *this = ToPaste{} }
func (this *ToPaste) String() string { return proto1.CompactTextString(this) }
func (*ToPaste) ProtoMessage()       {}

func (this *ToPaste) GetName() string {
	if this != nil && this.Name != nil {
		return *this.Name
	}
	return ""
}

func (this *ToPaste) GetData() []byte {
	if this != nil {
		return this.Data
	}
	return nil
}

type Posted struct {
	Url              *string `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *Posted) Reset()         { *this = Posted{} }
func (this *Posted) String() string { return proto1.CompactTextString(this) }
func (*Posted) ProtoMessage()       {}

func (this *Posted) GetUrl() string {
	if this != nil && this.Url != nil {
		return *this.Url
	}
	return ""
}

type Empty struct {
	XXX_unrecognized []byte `json:"-"`
}

func (this *Empty) Reset()         { *this = Empty{} }
func (this *Empty) String() string { return proto1.CompactTextString(this) }
func (*Empty) ProtoMessage()       {}

func init() {
}

// GoPaste is an interface satisfied by the generated client and
// which must be implemented by the object wrapped by the server.
type GoPaste interface {
	Paste(in *ToPaste, out *Posted) error
	Next(in *Empty, out *Posted) error
}

// internal wrapper for type-safe RPC calling
type rpcGoPasteClient struct {
	*rpc.Client
}

func (this rpcGoPasteClient) Paste(in *ToPaste, out *Posted) error {
	return this.Call("GoPaste.Paste", in, out)
}
func (this rpcGoPasteClient) Next(in *Empty, out *Posted) error {
	return this.Call("GoPaste.Next", in, out)
}

// NewGoPasteClient returns an *rpc.Client wrapper for calling the methods of
// GoPaste remotely.
func NewGoPasteClient(conn net.Conn) GoPaste {
	return rpcGoPasteClient{rpc.NewClientWithCodec(codec.NewClientCodec(conn))}
}

// ServeGoPaste serves the given GoPaste backend implementation on conn.
func ServeGoPaste(conn net.Conn, backend GoPaste) error {
	srv := rpc.NewServer()
	if err := srv.RegisterName("GoPaste", backend); err != nil {
		return err
	}
	srv.ServeCodec(codec.NewServerCodec(conn))
	return nil
}

// DialGoPaste returns a GoPaste for calling the GoPaste servince at addr (TCP).
func DialGoPaste(addr string) (GoPaste, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewGoPasteClient(conn), nil
}

// ListenAndServeGoPaste serves the given GoPaste backend implementation
// on all connections accepted as a result of listening on addr (TCP).
func ListenAndServeGoPaste(addr string, backend GoPaste) error {
	clients, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := rpc.NewServer()
	if err := srv.RegisterName("GoPaste", backend); err != nil {
		return err
	}
	for {
		conn, err := clients.Accept()
		if err != nil {
			return err
		}
		go srv.ServeCodec(codec.NewServerCodec(conn))
	}
	panic("unreachable")
}

// GoPasteWeb is the web-based RPC version of the interface which
// must be implemented by the object wrapped by the webrpc server.
type GoPasteWeb interface {
	Paste(r *http.Request, in *ToPaste, out *Posted) error
	Next(r *http.Request, in *Empty, out *Posted) error
}

// internal wrapper for type-safe webrpc calling
type rpcGoPasteWebClient struct {
	remote   *url.URL
	protocol webrpc.Protocol
}

func (this rpcGoPasteWebClient) Paste(in *ToPaste, out *Posted) error {
	return webrpc.Post(this.protocol, this.remote, "/GoPaste/Paste", in, out)
}

func (this rpcGoPasteWebClient) Next(in *Empty, out *Posted) error {
	return webrpc.Post(this.protocol, this.remote, "/GoPaste/Next", in, out)
}

// Register a GoPasteWeb implementation with the given webrpc ServeMux.
// If mux is nil, the default webrpc.ServeMux is used.
func RegisterGoPasteWeb(this GoPasteWeb, mux webrpc.ServeMux) error {
	if mux == nil {
		mux = webrpc.DefaultServeMux
	}
	if err := mux.Handle("/GoPaste/Paste", func(c *webrpc.Call) error {
		in, out := new(ToPaste), new(Posted)
		if err := c.ReadRequest(in); err != nil {
			return err
		}
		if err := this.Paste(c.Request, in, out); err != nil {
			return err
		}
		return c.WriteResponse(out)
	}); err != nil {
		return err
	}
	if err := mux.Handle("/GoPaste/Next", func(c *webrpc.Call) error {
		in, out := new(Empty), new(Posted)
		if err := c.ReadRequest(in); err != nil {
			return err
		}
		if err := this.Next(c.Request, in, out); err != nil {
			return err
		}
		return c.WriteResponse(out)
	}); err != nil {
		return err
	}
	return nil
}

// NewGoPasteWebClient returns a webrpc wrapper for calling the methods of GoPaste
// remotely via the web.  The remote URL is the base URL of the webrpc server.
func NewGoPasteWebClient(pro webrpc.Protocol, remote *url.URL) GoPaste {
	return rpcGoPasteWebClient{remote, pro}
}

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/logprefix"
	"github.com/EndCrystal/Server/network"
	"github.com/EndCrystal/Server/packet"
	plug "github.com/EndCrystal/Server/plugin"
	"nhooyr.io/websocket"
)

var PluginId string = "core:network:websocket"

func PluginMain(i plug.PluginInterface) error {
	i.RegisterNetworkProtocol("ws", creator)
	i.RegisterNetworkProtocol("ws+unix", creator)
	return nil
}

var ESchemeError = errors.New("Invalied scheme")

type Server struct {
	source  *http.Server
	fetcher chan network.ClientInstance
}

type Client struct {
	source  *websocket.Conn
	fetcher chan packet.Packet
	cancel  func()
}

func (c Client) SendPacket(pkt packet.Packet) (err error) {
	writter, err := c.source.Writer(context.TODO(), websocket.MessageBinary)
	if err != nil {
		return err
	}
	defer func() {
		err = writter.Close()
		if e := recover(); e != nil {
			var ok bool
			if err, ok = e.(error); !ok {
				err = fmt.Errorf("Unknown error %v", e)
			}
		}
	}()
	out := packed.MakeOutput(writter)
	packet.BuildPacket(pkt, out)
	return
}
func (c Client) GetFetcher() <-chan packet.Packet { return c.fetcher }
func (c Client) Disconnect()                      { c.cancel() }

func (s Server) Stop() {
	close(s.fetcher)
	s.source.Close()
}

func (s Server) GetFetcher() <-chan network.ClientInstance { return s.fetcher }

var opts = &websocket.AcceptOptions{
	Subprotocols:       []string{"endcrystal"},
	InsecureSkipVerify: true,
}

type privdata struct{}

func handler(res http.ResponseWriter, req *http.Request) {
	log := logprefix.Get("[websocket plugin] ")
	var c *websocket.Conn
	var err error
	c, err = websocket.Accept(res, req, opts)
	if err != nil {
		res.WriteHeader(500)
		return
	}
	defer c.Close(websocket.StatusInternalError, "fallthrough")
	ctx, cancel := context.WithCancel(req.Context())
	ch := ctx.Value(privdata{}).(chan network.ClientInstance)
	pktch := make(chan packet.Packet)
	defer close(pktch)
	ch <- Client{c, pktch, cancel}
	for {
		typ, reader, err := c.Reader(ctx)
		if err != nil {
			var ce *websocket.CloseError
			if errors.As(err, &ce) {
				break
			}
			log.Printf("Failed to read from ws: %v", err)
			c.Close(websocket.StatusInternalError, "failed to read")
			return
		}
		if typ != websocket.MessageBinary {
			c.Close(websocket.StatusUnsupportedData, "")
			return
		}
		in := packed.MakeInput(reader)
		var pkt packet.Packet
		pkt, err = packet.ParsePacket(in, packet.ClientSide, ^uint16(0))
		if err != nil {
			log.Print(err)
			c.Close(websocket.StatusProtocolError, "failed to parse")
			return
		}
		pktch <- pkt
	}
	c.Close(websocket.StatusNormalClosure, "")
}

func creator(u *url.URL) (network.Server, error) {
	log := logprefix.Get("[websocket plugin] ")
	server := new(http.Server)
	var listener net.Listener
	var err error
	var usePath bool = false
	switch u.Scheme {
	case "ws+unix":
		rpath := filepath.Join(u.Hostname(), u.EscapedPath())
		os.Remove(rpath)
		log.Printf("Listen websocket over unix socket: %s", rpath)
		listener, err = net.Listen("unix", rpath)
	case "ws":
		log.Printf("Listen websocket over tcp: %s", u.Host)
		listener, err = net.Listen("tcp", u.Host)
		usePath = true
	default:
		return nil, ESchemeError
	}
	if err != nil {
		return nil, err
	}
	if usePath {
		mux := http.NewServeMux()
		mux.HandleFunc(u.RequestURI(), handler)
		server.Handler = mux
	} else {
		server.Handler = http.HandlerFunc(handler)
	}
	ch := make(chan network.ClientInstance)
	ctx := context.WithValue(context.Background(), privdata{}, ch)
	server.BaseContext = func(net.Listener) context.Context { return ctx }
	go server.Serve(listener)
	return Server{server, ch}, nil
}

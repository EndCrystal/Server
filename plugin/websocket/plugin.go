package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/logprefix"
	"github.com/EndCrystal/Server/network"
	"github.com/EndCrystal/Server/packet"
	plug "github.com/EndCrystal/Server/plugin"
	"nhooyr.io/websocket"
)

// PluginID plugin id
var PluginID string = "core:network:websocket"

// Dependencies plugin dependencies
var Dependencies = []string{}

// PluginMain plugin main
func PluginMain(i plug.PluginInterface) error {
	i.RegisterNetworkProtocol("ws", creator)
	i.RegisterNetworkProtocol("ws+unix", creator)
	return nil
}

// ErrInvalidScheme Invalied scheme
var ErrInvalidScheme = errors.New("Invalied scheme")

// Server websocket server
type Server struct {
	source  *http.Server
	fetcher chan network.ClientInstance
}

// Client websocket client
type Client struct {
	source   *websocket.Conn
	identify network.CommonIdentifier
	fetcher  chan packet.ReceiveOnlyPacket
	cancel   func()
	mtx      *sync.Mutex
}

var bufpool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// SendPacket send packet
func (c Client) SendPacket(pkt packet.SendOnlyPacket) (err error) {
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			if err, ok = e.(error); !ok {
				err = fmt.Errorf("Unknown error %v", e)
			}
		}
	}()
	buf := bufpool.Get().(*bytes.Buffer)
	defer bufpool.Put(buf)
	buf.Reset()
	o := packed.MakeOutput(buf)
	packet.WritePacket(pkt, o)
	err = c.source.Write(context.TODO(), websocket.MessageBinary, buf.Bytes())
	return
}

// GetFetcher get fetch for packet
func (c Client) GetFetcher() <-chan packet.ReceiveOnlyPacket { return c.fetcher }

// Disconnect kick client
func (c Client) Disconnect() { c.cancel() }

// GetIdentifier get network identifier
func (c Client) GetIdentifier() network.Identifier { return c.identify }

// Stop stop the server
func (s Server) Stop() {
	close(s.fetcher)
	s.source.Close()
}

// GetFetcher get client fetcher
func (s Server) GetFetcher() <-chan network.ClientInstance { return s.fetcher }

var opts = &websocket.AcceptOptions{
	Subprotocols:       []string{"endcrystal"},
	InsecureSkipVerify: true,
}

type privdata struct{}

func getCommonNetworkIdentifier(req *http.Request) (id network.CommonIdentifier) {
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if len(ip) != 0 {
		id.IP = net.ParseIP(ip)
		return
	}
	ip = strings.TrimSpace(req.Header.Get("X-Real-Ip"))
	if len(ip) != 0 {
		id.IP = net.ParseIP(ip)
		return
	}
	if ip, port, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
		id.IP = net.ParseIP(ip)
		xport, err := strconv.Atoi(port)
		if err != nil {
			return
		}
		id.Port = uint16(xport)
		return
	}
	return
}

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
	pktch := make(chan packet.ReceiveOnlyPacket)
	defer close(pktch)
	ch <- Client{
		source:   c,
		identify: getCommonNetworkIdentifier(req),
		fetcher:  pktch,
		cancel:   cancel,
		mtx:      new(sync.Mutex),
	}
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
		var pkt packet.ReceiveOnlyPacket
		pkt, err = packet.ParsePacket(in, ^uint16(0))
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
		return nil, ErrInvalidScheme
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

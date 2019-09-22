package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/logprefix"
	"github.com/EndCrystal/Server/packet"
	"github.com/EndCrystal/Server/token"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Simple Agent")
}

var generator token.TokenGenerator

func auth(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	server := ps.ByName("server")
	username := ps.ByName("username")
	var payload packet.LoginPayload
	payload.ServerId = server
	payload.Username = username
	payload.Time = time.Now()
	var pkt packet.LoginPacket
	pkt.Write(payload, generator)

	out := packed.MakeOutput(w)
	pkt.Save(out)
}

var keypath = flag.String("key", "key.priv", "Path to server privkey")
var listen = flag.String("listen", ":1984", "Listen host")

func loadGenerator() {
	log := logprefix.Get("[privkey loader] ")
	log.Printf("Loading from %s", *keypath)
	var priv token.PrivKey
	generator = token.GetTokenGenerator(priv)
}

func main() {
	log := logprefix.Get("[simple auth] ")
	flag.Parse()

	loadGenerator()

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/login/:server/:username", auth)

	log.Print("Started")

	log.Fatal(http.ListenAndServe(*listen, router))
}

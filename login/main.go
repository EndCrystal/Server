package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/logprefix"
	"github.com/EndCrystal/Server/packet"
	"github.com/EndCrystal/Server/token"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Simple Agent")
}

var generator token.TokenGenerator

func genpacket(server, username string) (pkt packet.LoginPacket) {
	var payload packet.LoginPayload
	payload.ServerId = server
	payload.Username = username
	payload.Time = time.Now()
	pkt.Write(payload, generator)
	return
}

var log = logprefix.Get("[login server] ")

func auth(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	server := ps.ByName("server")
	username := ps.ByName("username")

	go log.Printf("Generate %s for %s", username, server)
	pkt := genpacket(server, username)

	out := packed.MakeOutput(w)
	packet.WritePacket(&pkt, out)
}

var keypath = flag.String("key", "key.priv", "Path to server privkey")
var listen = flag.String("listen", ":1984", "Listen host")

func loadGenerator() {
	log := logprefix.Get("[privkey loader] ")
	log.Printf("Loading from %s", *keypath)
	stat, err := os.Stat(*keypath)
	if err != nil {
		return
	}
	if stat.Size() != int64(token.PrivKeyLen) {
		panic("Failed to load privkey: size mismatch")
	}
	data, err := ioutil.ReadFile(*keypath)
	if err != nil {
		panic(err)
	}
	var priv token.PrivKey
	copy(priv[:], data)
	generator = token.GetTokenGenerator(priv)
}

func main() {
	log := logprefix.Get("[simple auth] ")
	flag.Parse()

	loadGenerator()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/login/:server/:username", auth)

	handler := cors.Default().Handler(router)

	log.Printf("Started %s", *listen)

	log.Fatal(http.ListenAndServe(*listen, handler))
}

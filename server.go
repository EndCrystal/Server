package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/EndCrystal/Server/logprefix"
	"github.com/EndCrystal/Server/network"
	plug "github.com/EndCrystal/Server/plugin"
	"github.com/EndCrystal/Server/token"
)

type ChatMessage struct {
	Sender  string
	Message string
}

var global struct {
	verfier token.TokenVerifier
	chat    chan<- ChatMessage
}

var log = logprefix.Get("[main] ")

func main() {
	var err error
	flag.Parse()
	err = loadPluginFromMulti(strings.Split(*plugin_home, ":")...)
	if err != nil {
		log.Fatalf("Failed to load plugins: %v", err)
	}
	printLoadedPlugins()

	global.verfier, err = loadPubKey()
	if err != nil {
		log.Fatalf("Failed to load pubkey: %v", err)
	}
	// TODO: global.chat

	var server network.Server
	var endpoint_url *url.URL
	endpoint_url, err = url.Parse(*endpoint)
	if err != nil {
		log.Fatalf("Failed to parse endpoint url: %v", err)
	}
	server, err = network.CreateServer(endpoint_url)
	if err != nil {
		log.Fatalf("Failed to create server for this endpoint (%s): %v", *endpoint, err)
	}
	defer server.Stop()
	loop(server.GetFetcher())
}

var endpoint = flag.String("endpoint", "ws://0.0.0.0:2480", "Server Endpoint")
var plugin_home = flag.String("plugin-dirs", "plugins:"+filepath.Join(os.Getenv("HOME"), ".local", "share", "EndCrystal", "plugins"), "Plugin directories")
var pubkey_path = flag.String("pubkey", "key.pub", "Path to server pubkey")
var server_id = flag.String("server-id", "default", "Server Id")
var connection_timeout = flag.Duration("verify-timeout", time.Second*10, "Timeout for verify login packet")

func loadPubKey() (verifier token.TokenVerifier, err error) {
	log := logprefix.Get("[pubkey loader] ")
	log.Printf("Loading from %s", *pubkey_path)
	stat, err := os.Stat(*pubkey_path)
	if err != nil {
		return
	}
	if stat.Size() != int64(token.PubKeyLen) {
		return nil, fmt.Errorf("Failed to load pubkey: size mismatch")
	}
	data, err := ioutil.ReadFile(*pubkey_path)
	if err != nil {
		return
	}
	var pubkey token.PubKey
	copy(pubkey[:], data)
	verifier = token.GetTokenVerifier(pubkey)
	log.Printf("Loaded")
	return
}

func printLoadedPlugins() {
	for id := range plug.LoadedPlugins {
		log.Printf("Loaded plugin: %s", id)
	}
}

func loadPluginFromMulti(roots ...string) (err error) {
	for _, root := range roots {
		err = loadPluginFrom(root)
		if err != nil {
			return
		}
	}
	return
}

func loadPluginFrom(root string) (err error) {
	log := logprefix.Get("[plugin loader] ")
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	var plugin_count int
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".ecplugin") {
			return nil
		}
		if info.Mode().Perm()&0111 != 0111 {
			log.Printf("Skiped plugin: %s", path)
			return nil
		}
		log.Printf("Loading plugin: %s", path)
		err = plug.LoadPlugin(path, plug.PluginInterface{})
		if err != nil {
			return err
		}
		plugin_count++
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Skipped loading plugins from %s", root)
			err = nil
			return
		}
		return
	}
	log.Printf("Loaded %d plugins from %s", plugin_count, root)
	return
}

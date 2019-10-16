package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/EndCrystal/Server/logprefix"
	"github.com/EndCrystal/Server/network"
	plug "github.com/EndCrystal/Server/plugin"
	"github.com/EndCrystal/Server/token"
	"github.com/EndCrystal/Server/world/chunk"
	"github.com/EndCrystal/Server/world/storage"
	go_up "github.com/ufoscout/go-up"
)

type chatMessage struct {
	Sender  string
	Message string
}

type serverGlobal struct {
	verfier token.Verifier
	chat    chan<- chatMessage
	users   *sync.Map
}

var global serverGlobal

var log = logprefix.Get("[main] ")

func main() {
	var err error
	loadConfig()
	err = loadMainStorage()
	if err != nil {
		log.Fatalf("Failed to load main storage: %s", config.storagePath)
	}
	defer storage.MainStorage.Close()
	err = loadPluginFromMulti(config.pluginHome...)
	if err != nil {
		pluginStats()
		log.Fatalf("Failed to load plugins: %v", err)
	}
	log.Print("Applying plugin...")
	err = plug.ApplyPlugin(plug.PluginInterface{})
	if err != nil {
		pluginStats()
		log.Fatalf("Failed to apply plugin: %v", err)
	}
	pluginStats()

	global.verfier, err = loadPubKey()
	if err != nil {
		log.Fatalf("Failed to load pubkey: %v", err)
	}
	global.users = new(sync.Map)
	global.chat = handleChat(global)

	var server network.Server
	var endpointURL *url.URL
	endpointURL, err = url.Parse(config.endpoint)
	if err != nil {
		log.Fatalf("Failed to parse endpoint url: %v", err)
	}
	server, err = network.CreateServer(endpointURL)
	if err != nil {
		log.Fatalf("Failed to create server for this endpoint (%s): %v", config.endpoint, err)
	}
	defer server.Stop()
	done := make(chan struct{})
	go loop(server.GetFetcher(), done)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-done:
	case sig := <-sigs:
		log.Printf("Received signal %v: exiting...", sig)
	}
}

func loadConfig() {
	log := logprefix.Get("[config loader] ")
	var err error
	checkerr := func() {
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	}
	up, err := go_up.NewGoUp().
		AddFile("/etc/EndCrystal/config.properties", true).
		AddFile(filepath.Join(os.Getenv("HOME"), ".config", "EndCrystal", "config.properties"), true).
		AddFile("config.properties", true).
		AddReader(go_up.NewEnvReader("EC_", true, true)).
		Build()
	checkerr()
	config.endpoint = up.GetStringOrDefault("endpoint", "ws://0.0.0.0:2480")
	config.pubkeyPath = up.GetStringOrDefault("pubkey.path", "key.pub")
	config.pluginHome = up.GetStringSliceOrDefault("plugin", ":", []string{"plugins", filepath.Join(os.Getenv("HOME"), ".local", "share", "EndCrystal", "plugins")})
	config.id = up.GetStringOrDefault("id", "default")
	config.connectionTimeout = time.Duration(up.GetIntOrDefault("connection.timeout", 10)) * time.Second
	config.storagePath = up.GetStringOrDefault("storage", "EndCrystal.bbolt")
	config.spawnpoint.dimension, err = up.GetStringOrFail("spawnpoint.dimension")
	checkerr()
	config.viewDistance = uint32(up.GetIntOrDefault("view-distance", 8))

	{
		var pos string
		pos, err = up.GetStringOrFail("spawnpoint.pos")
		fmt.Sscanf(pos, "%d:%d", &config.spawnpoint.pos.X, &config.spawnpoint.pos.Z)
	}
}

var config struct {
	endpoint          string
	pluginHome        []string
	pubkeyPath        string
	id                string
	connectionTimeout time.Duration
	storagePath       string
	viewDistance      uint32
	spawnpoint        struct {
		dimension string
		pos       chunk.CPos
	}
}

func loadMainStorage() (err error) {
	storage.MainStorage, err = storage.Open(config.storagePath)
	return
}

func loadPubKey() (verifier token.Verifier, err error) {
	log := logprefix.Get("[pubkey loader] ")
	log.Printf("Loading from %s", config.pubkeyPath)
	stat, err := os.Stat(config.pubkeyPath)
	if err != nil {
		return
	}
	if stat.Size() != int64(token.PubKeyLen) {
		return nil, fmt.Errorf("Failed to load pubkey: size mismatch")
	}
	data, err := ioutil.ReadFile(config.pubkeyPath)
	if err != nil {
		return
	}
	var pubkey token.PubKey
	copy(pubkey[:], data)
	verifier = token.GetTokenVerifier(pubkey)
	return
}

func pluginStats() {
	log := logprefix.Get("[plugin stats] ")
	for id := range plug.PendingPlugins {
		log.Printf("Queued %s", id)
	}
	for id := range plug.LoadedPlugins {
		log.Printf("Loaded %s", id)
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
	var pluginCount int
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
		log.Printf("Found plugin: %s", path)
		err = plug.LoadPlugin(path)
		if err != nil {
			return err
		}
		pluginCount++
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
	log.Printf("Queued %d plugins from %s", pluginCount, root)
	return
}

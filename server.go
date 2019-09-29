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

type ChatMessage struct {
	Sender  string
	Message string
}

type Global struct {
	verfier token.TokenVerifier
	chat    chan<- ChatMessage
	users   *sync.Map
}

var global Global

var log = logprefix.Get("[main] ")

func main() {
	var err error
	loadConfig()
	err = loadMainStorage()
	if err != nil {
		log.Fatalf("Failed to load main storage: %s", config.storage_path)
	}
	defer storage.MainStorage.Close()
	err = loadPluginFromMulti(config.plugin_home...)
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
	var endpoint_url *url.URL
	endpoint_url, err = url.Parse(config.endpoint)
	if err != nil {
		log.Fatalf("Failed to parse endpoint url: %v", err)
	}
	server, err = network.CreateServer(endpoint_url)
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
	config.pubkey_path = up.GetStringOrDefault("pubkey.path", "key.pub")
	config.plugin_home = up.GetStringSliceOrDefault("plugin", ":", []string{"plugins", filepath.Join(os.Getenv("HOME"), ".local", "share", "EndCrystal", "plugins")})
	config.id = up.GetStringOrDefault("id", "default")
	config.connection_timeout = time.Duration(up.GetIntOrDefault("connection.timeout", 10)) * time.Second
	config.storage_path = up.GetStringOrDefault("storage", "EndCrystal.bbolt")
	config.spawnpoint.dimension, err = up.GetStringOrFail("spawnpoint.dimension")
	checkerr()

	{
		var pos string
		pos, err = up.GetStringOrFail("spawnpoint.pos")
		fmt.Sscanf(pos, "%d:%d", &config.spawnpoint.pos.X, &config.spawnpoint.pos.Z)
	}
}

var config struct {
	endpoint           string
	plugin_home        []string
	pubkey_path        string
	id                 string
	connection_timeout time.Duration
	storage_path       string
	spawnpoint         struct {
		dimension string
		pos       chunk.ChunkPos
	}
}

func loadMainStorage() (err error) {
	storage.MainStorage, err = storage.Open(config.storage_path)
	return
}

func loadPubKey() (verifier token.TokenVerifier, err error) {
	log := logprefix.Get("[pubkey loader] ")
	log.Printf("Loading from %s", config.pubkey_path)
	stat, err := os.Stat(config.pubkey_path)
	if err != nil {
		return
	}
	if stat.Size() != int64(token.PubKeyLen) {
		return nil, fmt.Errorf("Failed to load pubkey: size mismatch")
	}
	data, err := ioutil.ReadFile(config.pubkey_path)
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
		log.Printf("Found plugin: %s", path)
		err = plug.LoadPlugin(path)
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
	log.Printf("Queued %d plugins from %s", plugin_count, root)
	return
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	. "github.com/EndCrystal/Server/logprefix"
	plug "github.com/EndCrystal/Server/plugin"
)

func main() {
	defer LogPrefix(LogPrefix("[main] "))
	var err error
	flag.Parse()
	err = loadPluginFromMulti(append(strings.Split(*plugin_home, ":"), filepath.Join(os.Getenv("HOME"), ".local", "share", "EndCrystal", "plugins"))...)
	if err != nil {
		log.Fatalf("Failed to load plugins: %v", err)
	}
	printLoadedPlugins()
}

var endpoint = flag.String("endpoint", "ws://0.0.0.0:2480", "Server Endpoint")
var plugin_home = flag.String("plugin-dirs", "plugins", "Plugin directories")

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

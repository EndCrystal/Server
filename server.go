package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	plug "github.com/EndCrystal/Server/plugin"
)

func main() {
	var err error
	err = loadPluginFrom("plugins")
	if err != nil {
		log.Panicf("Failed to load plugins: %v", err)
	}
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
	err = filepath.Walk("plugins", func(path string, info os.FileInfo, err error) error {
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
		return
	}
	log.Printf("Loaded %d plugins", plugin_count)
	return
}

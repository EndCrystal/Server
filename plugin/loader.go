package plug

import (
	"errors"
	"fmt"
	"plugin"
)

var EIncompatiablePlugin = errors.New("Incomparable plugin")
var EConflictPlugin = errors.New("Conflict plugin")
var EInvalidPluginId = errors.New("Invalid plugin id")
var EInvalidResolver = errors.New("Invalid plugin resolver")

type EInvalidPlugin struct{ wrapped error }

func (e EInvalidPlugin) Unwrap() error { return e.wrapped }
func (e EInvalidPlugin) Error() string { return fmt.Sprintf("Invalid plugin: %v", e.wrapped) }

var LoadedPlugins map[string]*plugin.Plugin

const PluginSystemVersion = 0

// Load plugin from path
// plugin sample:
// var PluginId string = "My test mod"
// func PluginMain(plug.PluginInterface) error
func LoadPlugin(path string, pifce PluginInterface) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}
	sname, err := p.Lookup("PluginId")
	if err != nil {
		return EInvalidPlugin{err}
	}
	name, ok := sname.(string)
	if !ok {
		return EInvalidPlugin{EInvalidPluginId}
	}
	if _, ok := LoadedPlugins[name]; ok {
		return EConflictPlugin
	}
	fn, err := loadPluginMain(p, "PluginMain", "PluginMainV0")
	err = fn(pifce)
	if err != nil {
		return EInvalidPlugin{err}
	}
	LoadedPlugins[name] = p
	return nil
}

func loadPluginMain(p *plugin.Plugin, arr ...string) (func(PluginInterface) error, error) {
	for _, mainName := range arr {
		sym, err := p.Lookup(mainName)
		if err != nil {
			return nil, EInvalidPlugin{err}
		}
		fn, ok := sym.(func(PluginInterface) error)
		if ok {
			return fn, nil
		}
	}
	return nil, EIncompatiablePlugin
}

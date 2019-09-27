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

type EBrokenDependencies struct{ Name, Source string }

func (e EBrokenDependencies) Error() string {
	return fmt.Sprintf("Broken dependency: %s not found (required by %s)", e.Name, e.Source)
}

type EInvalidPlugin struct{ wrapped error }

func (e EInvalidPlugin) Unwrap() error { return e.wrapped }
func (e EInvalidPlugin) Error() string { return fmt.Sprintf("Invalid plugin: %v", e.wrapped) }

type PluginInfo struct {
	*plugin.Plugin
	dependencies []string
	fn           func(PluginInterface) error
}

var PendingPlugins = make(map[string]*PluginInfo)
var LoadedPlugins = make(map[string]*PluginInfo)

const PluginSystemVersion = 0

// Load plugin from path
// plugin sample:
// var PluginId string = "My test mod"
// var Dependencies []string
// func PluginMain(plug.PluginInterface) error
func LoadPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}
	sname, err := p.Lookup("PluginId")
	if err != nil {
		return EInvalidPlugin{err}
	}
	pname, ok := sname.(*string)
	if !ok {
		return EInvalidPlugin{EInvalidPluginId}
	}
	name := *pname
	if _, ok := LoadedPlugins[name]; ok {
		return EConflictPlugin
	}
	deps, err := loadPluginDependencies(p)
	if err != nil {
		return EInvalidPlugin{err}
	}
	fn, err := loadPluginMain(p, "PluginMain", "PluginMainV0")
	if err != nil {
		return EInvalidPlugin{err}
	}
	PendingPlugins[name] = &PluginInfo{Plugin: p, dependencies: deps, fn: fn}
	return nil
}

func sortPlugins(ch chan<- *PluginInfo) error {
	defer close(ch)
	seen := make(map[string]bool)
	var visitAll func(source string, items []string) error
	visitAll = func(source string, items []string) error {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				if pp, ok := PendingPlugins[item]; ok {
					err := visitAll(item, pp.dependencies)
					if err != nil {
						return err
					}
					ch <- pp
					LoadedPlugins[item] = pp
					delete(PendingPlugins, item)
				} else {
					return EBrokenDependencies{item, source}
				}
			}
		}
		return nil
	}
	rootdeps := make([]string, len(PendingPlugins))
	i := 0
	for name := range PendingPlugins {
		rootdeps[i] = name
		i++
	}
	return visitAll("", rootdeps)
}

func ApplyPlugin(pifce PluginInterface) error {
	ch := make(chan *PluginInfo)
	errch := make(chan error)
	go func() {
		if err := sortPlugins(ch); err != nil {
			errch <- err
		}
		close(errch)
	}()
	for {
		select {
		case err := <-errch:
			return err
		case item, ok := <-ch:
			if !ok {
				return nil
			}
			if err := item.fn(pifce); err != nil {
				return err
			}
		}
	}
}

func loadPluginDependencies(p *plugin.Plugin) ([]string, error) {
	sym, err := p.Lookup("Dependencies")
	if err != nil {
		return nil, EInvalidPlugin{err}
	}
	ret, ok := sym.(*[]string)
	if ok {
		return *ret, nil
	}
	return nil, EIncompatiablePlugin
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

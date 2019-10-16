package plug

import (
	"errors"
	"fmt"
	"plugin"
)

// ErrIncompatiablePlugin Incomparable plugin
var ErrIncompatiablePlugin = errors.New("Incomparable plugin")

// ErrConflictPlugin Conflict plugin
var ErrConflictPlugin = errors.New("Conflict plugin")

// ErrInvalidPluginID Invalid plugin id
var ErrInvalidPluginID = errors.New("Invalid plugin id")

// ErrInvalidResolver Invalid plugin resolver
var ErrInvalidResolver = errors.New("Invalid plugin resolver")

// ErrBrokenDependencies broken dependency
type ErrBrokenDependencies struct{ Name, Source string }

func (e ErrBrokenDependencies) Error() string {
	return fmt.Sprintf("Broken dependency: %s not found (required by %s)", e.Name, e.Source)
}

// ErrInvalidPlugin Invalid plugin
type ErrInvalidPlugin struct{ wrapped error }

// Unwrap unwrap for nested error
func (e ErrInvalidPlugin) Unwrap() error { return e.wrapped }
func (e ErrInvalidPlugin) Error() string { return fmt.Sprintf("Invalid plugin: %v", e.wrapped) }

// PluginInfo basic plugin info
type PluginInfo struct {
	*plugin.Plugin
	dependencies []string
	fn           func(PluginInterface) error
}

// PendingPlugins plugins that pending to load
var PendingPlugins = make(map[string]*PluginInfo)

// LoadedPlugins loaded plugins
var LoadedPlugins = make(map[string]*PluginInfo)

// PluginSystemVersion plugin system version
const PluginSystemVersion = 0

// LoadPlugin Load plugin from path
// plugin sample:
// var PluginID string = "My test mod"
// var Dependencies []string
// func PluginMain(plug.PluginInterface) error
func LoadPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}
	sname, err := p.Lookup("PluginID")
	if err != nil {
		return ErrInvalidPlugin{err}
	}
	pname, ok := sname.(*string)
	if !ok {
		return ErrInvalidPlugin{ErrInvalidPluginID}
	}
	name := *pname
	if _, ok := LoadedPlugins[name]; ok {
		return ErrConflictPlugin
	}
	deps, err := loadPluginDependencies(p)
	if err != nil {
		return ErrInvalidPlugin{err}
	}
	fn, err := loadPluginMain(p, "PluginMain", "PluginMainV0")
	if err != nil {
		return ErrInvalidPlugin{err}
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
					return ErrBrokenDependencies{item, source}
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

// ApplyPlugin apply plugins
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
		return nil, ErrInvalidPlugin{err}
	}
	ret, ok := sym.(*[]string)
	if ok {
		return *ret, nil
	}
	return nil, ErrIncompatiablePlugin
}

func loadPluginMain(p *plugin.Plugin, arr ...string) (func(PluginInterface) error, error) {
	for _, mainName := range arr {
		sym, err := p.Lookup(mainName)
		if err != nil {
			return nil, ErrInvalidPlugin{err}
		}
		fn, ok := sym.(func(PluginInterface) error)
		if ok {
			return fn, nil
		}
	}
	return nil, ErrIncompatiablePlugin
}

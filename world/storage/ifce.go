package storage

// MainStorage main storage
var MainStorage *Storage

// PluginStorageHost plugin host
type PluginStorageHost struct{}

// GetMainStorage get main storage
func (PluginStorageHost) GetMainStorage() *Storage {
	return MainStorage
}

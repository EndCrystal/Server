package storage

var MainStorage *Storage

type PluginStorageHost struct{}

func (PluginStorageHost) GetMainStorage() *Storage {
	return MainStorage
}

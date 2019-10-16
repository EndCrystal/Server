package block

// PluginBlockHost plugin host
type PluginBlockHost struct{}

// RegisterBlock register block
func (PluginBlockHost) RegisterBlock(b Block) BID {
	return Register(b)
}

// GetBlock get block by id
func (PluginBlockHost) GetBlock(id BID) *Block {
	return Get(id)
}

// LookupBlockID lookup block id by name
func (PluginBlockHost) LookupBlockID(name string) (BID, bool) {
	return LookupID(name)
}

// LookupBlock lookup block by name
func (PluginBlockHost) LookupBlock(name string) *Block {
	return Lookup(name)
}

// GetBlockExtraLoader get block extra data loader
func (PluginBlockHost) GetBlockExtraLoader(id BID) ExtraLoader {
	return GetExtraLoader(id)
}

// RegisterBlockExtraLoader register block extra loader
func (PluginBlockHost) RegisterBlockExtraLoader(id BID, loader ExtraLoader) {
	RegisterExtraLoader(id, loader)
}

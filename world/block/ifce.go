package block

type PluginBlockHost struct{}

func (PluginBlockHost) RegisterBlock(b Block) Id {
	return Register(b)
}

func (PluginBlockHost) GetBlock(id Id) *Block {
	return Get(id)
}

func (PluginBlockHost) LookupBlockId(name string) (Id, bool) {
	return LookupId(name)
}

func (PluginBlockHost) LookupBlock(name string) *Block {
	return Lookup(name)
}

func (PluginBlockHost) GetBlockExtraLoader(id Id) ExtraLoader {
	return GetExtraLoader(id)
}

func (PluginBlockHost) RegisterBlockExtraLoader(id Id, loader ExtraLoader) {
	RegisterExtraLoader(id, loader)
}

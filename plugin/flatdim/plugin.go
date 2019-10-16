package main

import (
	packed "github.com/EndCrystal/PackedIO"
	plug "github.com/EndCrystal/Server/plugin"
	"github.com/EndCrystal/Server/world/chunk"
)

// PluginID plugin identifier
var PluginID string = "core:dim:flatdim"
// Dependencies dependencies
var Dependencies = []string{}

// PluginMain plugin entry
func PluginMain(i plug.PluginInterface) error {
	s := i.GetMainStorage().ForDim("flat")
	i.AddDimension("flat", []string{"flat", "test"}, s, SimpleFlatWorldGenerator{i})
	return nil
}

// SimpleFlatWorldGenerator simple flat world generator
type SimpleFlatWorldGenerator struct{ ifce plug.PluginInterface }

// Load load from data
func (SimpleFlatWorldGenerator) Load(packed.Input)  {}
// Save save to data
func (SimpleFlatWorldGenerator) Save(packed.Output) {}
// Generate generate chunk
func (g SimpleFlatWorldGenerator) Generate(pos chunk.CPos) *chunk.Chunk {
	ret := new(chunk.Chunk)
	bedrock, ok := g.ifce.LookupBlockID("core:bedrock")
	if !ok {
		panic("Cannot found core:bedrock")
	}
	for i := 0; i < 16*16*4; i++ {
		ret.Blocks[i].ID = bedrock
	}
	return ret
}

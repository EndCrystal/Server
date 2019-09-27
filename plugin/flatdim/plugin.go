package main

import (
	packed "github.com/EndCrystal/PackedIO"
	plug "github.com/EndCrystal/Server/plugin"
	"github.com/EndCrystal/Server/world/chunk"
)

var PluginId string = "core:dim:flatdim"
var Dependencies = []string{}

func PluginMain(i plug.PluginInterface) error {
	s := i.GetMainStorage().ForDim("flat")
	i.AddDimension("flat", []string{"flat", "test"}, s, SimpleFlatWorldGenerator{i})
	return nil
}

type SimpleFlatWorldGenerator struct{ ifce plug.PluginInterface }

func (SimpleFlatWorldGenerator) Load(packed.Input)  {}
func (SimpleFlatWorldGenerator) Save(packed.Output) {}
func (g SimpleFlatWorldGenerator) Generate(pos chunk.ChunkPos) *chunk.Chunk {
	ret := new(chunk.Chunk)
	bedrock, ok := g.ifce.LookupBlockId("core:bedrock")
	if !ok {
		panic("Cannot found core:bedrock")
	}
	for i := 0; i < 16*16*4; i++ {
		ret.Blocks[i].Id = bedrock
	}
	return ret
}

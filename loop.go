package main

import (
	"github.com/EndCrystal/Server/network"
	"github.com/EndCrystal/Server/packet"
)

func loop(ch <-chan network.ClientInstance) {
	for instance := range ch {
		go processClient(instance)
	}
}

type clientState struct {
	username string
}

func processClient(instance network.ClientInstance) {
	var state clientState
	var err error
	state, err = processLogin(instance)
	if err != nil {
		return
	}
	fetcher := instance.GetFetcher()
	for packet := range fetcher {
		processPacket(instance, &state, packet)
	}
}

func processLogin(instance network.ClientInstance) (state clientState, err error) {
	return
}

func processPacket(instance network.ClientInstance, state *clientState, pkt packet.Packet) {
}

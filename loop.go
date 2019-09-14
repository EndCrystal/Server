package main

import (
	"fmt"
	"time"

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
	defer func() {
		if err != nil {
			instance.SendPacket(&packet.DisconnectPacket{err.Error()})
		}
		instance.Disconnect()
	}()
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
	pkt, ok := <-instance.GetFetcher()
	if !ok {
		err = fmt.Errorf("failed to login")
		return
	}
	login_pkt, ok := pkt.(*packet.LoginPacket)
	if !ok {
		err = fmt.Errorf("failed to login: the first packet should be LoginPacket")
		return
	}
	if !login_pkt.Verify(global.verfier) {
		err = fmt.Errorf("Verify failed")
		return
	}
	payload, ok := login_pkt.Read()
	if !ok {
		err = fmt.Errorf("Cannot read payload")
		return
	}
	if payload.ServerId != *server_id {
		err = fmt.Errorf("Server id mismatch")
		return
	}
	if time.Since(payload.Time) > time.Second*10 {
		err = fmt.Errorf("Timeout")
		return
	}
	if len(payload.Username) == 0 {
		err = fmt.Errorf("Illegal username")
		return
	}
	state.username = payload.Username
	return
}

func processPacket(instance network.ClientInstance, state *clientState, pkt packet.Packet) {
}

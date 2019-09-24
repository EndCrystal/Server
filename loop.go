package main

import (
	"fmt"
	"time"

	"github.com/EndCrystal/Server/common"
	"github.com/EndCrystal/Server/logprefix"
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
	log := logprefix.Get("[process client] ")
	log.Println("New connection")
	defer log.Println("Client left")
	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
			if err, ok := r.(error); ok {
				instance.SendPacket(&packet.DisconnectPacket{Message: err.Error()})
			} else {
				instance.SendPacket(&packet.DisconnectPacket{Message: "unknown error"})
			}
		}
		instance.Disconnect()
	}()
	var state clientState
	var err error
	state, err = processLoginWithTimeout(instance, 5*time.Second)
	if err != nil {
		panic(err)
	}
	log.Printf("Player %s joined: %v", state.username, instance.GetIdentifier())
	global.users.Store(state.username, instance)
	defer global.users.Delete(state.username)
	startpacket := common.Value.GameStartHandler(state.username, instance.GetIdentifier())
	instance.SendPacket(&startpacket)
	fetcher := instance.GetFetcher()
	for packet := range fetcher {
		processPacket(instance, &state, packet)
	}
}

func processLoginWithTimeout(instance network.ClientInstance, timeout time.Duration) (state clientState, err error) {
	statech := make(chan clientState)
	errch := make(chan error)
	go func() {
		var state clientState
		var err error
		state, err = processLogin(instance)
		if err != nil {
			errch <- err
		}
		statech <- state
	}()
	timer := time.NewTimer(5 * time.Second)
	select {
	case state = <-statech:
		timer.Stop()
	case err = <-errch:
		timer.Stop()
	case <-timer.C:
		err = fmt.Errorf("timeout")
	}
	return
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
	if _, ok := global.users.Load(payload.Username); ok {
		err = fmt.Errorf("Forbidden login")
		return
	}
	state.username = payload.Username
	return
}

func processPacket(instance network.ClientInstance, state *clientState, pkt packet.Packet) {
	switch p := pkt.(type) {
	case *packet.LoginPacket:
		panic(fmt.Errorf("State mismatch"))
	case *packet.ChatPacket:
		msg := p.Message
		global.chat <- ChatMessage{state.username, msg}
	default:
		panic(fmt.Errorf("Unknown packet"))
	}
}

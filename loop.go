package main

import (
	"fmt"
	"time"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/common"
	"github.com/EndCrystal/Server/logprefix"
	"github.com/EndCrystal/Server/network"
	"github.com/EndCrystal/Server/packet"
	"github.com/EndCrystal/Server/user"
	"github.com/EndCrystal/Server/world/chunk"
	"github.com/EndCrystal/Server/world/dim"
	"github.com/EndCrystal/Server/world/storage"
)

func loop(ch <-chan network.ClientInstance, done chan<- struct{}) {
	for instance := range ch {
		go processClient(instance)
	}
	done <- struct{}{}
}

type userInfoWithClient struct {
	network.ClientInstance
	*user.Info
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
	var state user.Info
	var err error
	state, err = processLoginWithTimeout(instance, 5*time.Second)
	if err != nil {
		panic(err)
	}
	log.Printf("Player %s joined: %v", state.Username, instance.GetIdentifier())
	defer log.Printf("Player %s left: %v", state.Username, instance.GetIdentifier())
	global.users.Store(state.Username, userInfoWithClient{instance, &state})
	defer global.users.Delete(state.Username)

	var startpacket packet.GameStartPacket
	startpacket.Username = state.Username
	startpacket.Label = common.Value.UserLabelHandler(state.Username, instance.GetIdentifier())
	startpacket.Motd = common.Value.ServerMotdHandler(state.Username, instance.GetIdentifier())
	startpacket.MaxViewDistance = config.viewDistance
	startpacket.InitialPosition = state.Pos
	instance.SendPacket(&startpacket)

	fetcher := instance.GetFetcher()
	for packet := range fetcher {
		processPacket(instance, &state, packet)
	}
}

func processLoginWithTimeout(instance network.ClientInstance, timeout time.Duration) (state user.Info, err error) {
	statech := make(chan user.Info)
	errch := make(chan error)
	go func() {
		var state user.Info
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

func processLogin(instance network.ClientInstance) (state user.Info, err error) {
	pkt, ok := <-instance.GetFetcher()
	if !ok {
		err = fmt.Errorf("failed to login")
		return
	}
	loginPacket, ok := pkt.(*packet.LoginPacket)
	if !ok {
		err = fmt.Errorf("failed to login: the first packet should be LoginPacket")
		return
	}
	if !loginPacket.Verify(global.verfier) {
		err = fmt.Errorf("Verify failed")
		return
	}
	payload, ok := loginPacket.Read()
	if !ok {
		err = fmt.Errorf("Cannot read payload")
		return
	}
	if payload.ServerID != config.id {
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
	state.Username = payload.Username
	state.UserLabel = common.Value.UserLabelHandler(state.Username, instance.GetIdentifier())
	var dimensionName string
	var userConfig = storage.MainStorage.ForConfig("user")
	if data := userConfig.Get(state.Username); data != nil {
		// TODO: Just placeholder
		i := packed.InputFromBuffer(data)
		dimensionName = i.ReadString()
		state.Pos.Load(i)
	} else {
		dimensionName = config.spawnpoint.dimension
		state.Pos = config.spawnpoint.pos
	}
	if state.Dimension, ok = dim.LookupDimension(dimensionName); !ok {
		panic(fmt.Errorf("Failed to load dimension %v", dimensionName))
	}
	return
}

func processPacket(instance network.ClientInstance, state *user.Info, pkt packet.ReceiveOnlyPacket) {
	switch p := pkt.(type) {
	case *packet.BatchPacket:
		for _, sub := range p.ReceivedPackets {
			processPacket(instance, state, sub)
		}
	case *packet.LoginPacket:
		panic(fmt.Errorf("State mismatch"))
	case *packet.ChatPacket:
		msg := p.Message
		global.chat <- chatMessage{state.Username, msg}
	case *packet.ChunkRequestPacket:
		if state.Pos.Distance(p.Pos) <= config.viewDistance {
			fetchChunkForUser(instance, state, p.Pos)
		}
	default:
		panic(fmt.Errorf("Unknown packet %T", p))
	}
}

func fetchChunkForUser(instance network.ClientInstance, state *user.Info, pos chunk.CPos) {
	data, err := state.Dimension.GetChunk(pos)
	if err != nil {
		fmt.Printf("Failed to load chunk (@%+v) for user %s", pos, state.Username)
		return
	}
	defer data.Access()()
	pkt := packet.ChunkDataPacket{
		Pos: pos,
	}
	dest := pkt.SetData(data.Chunk)
	go func() {
		instance.SendPacket(pkt)
		dest()
	}()
}

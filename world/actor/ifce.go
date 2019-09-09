package actor

type PluginActorHost struct{}

func (PluginActorHost) AddActorSystem(sys System) { AddSystem(sys) }

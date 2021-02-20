package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ClientCommandRewrite struct{}

func (ClientCommandRewrite) Handle(_ Session, pk *packet.Packet) bool {
	switch pk2 := (*pk).(type) {
	case *packet.AvailableCommands:
		pk2.Commands = append(pk2.Commands, protocol.Command{
			Name:        "goproxy",
			Description: "GoProxy",
		})
	}
	return false
}

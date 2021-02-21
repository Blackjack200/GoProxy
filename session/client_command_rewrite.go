package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ClientCommandRewrite struct{}

func (ClientCommandRewrite) Handle(_ Session, pk *packet.Packet) bool {
	switch pk2 := (*pk).(type) {
	case *packet.AvailableCommands:
		commands := pk2.Commands

		for name := range Commands {
			commands = append(commands, protocol.Command{
				Name:        "__" + name,
				Description: "GoProxy Command",
			})
		}
		pk2.Commands = commands
	}
	return false
}

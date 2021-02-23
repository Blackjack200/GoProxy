package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ServerCommandRewrite struct{}

func (ServerCommandRewrite) Handle(_ *Session, pk *packet.Packet) bool {
	pk2, ok := (*pk).(*packet.AvailableCommands)
	if ok {
		commands := pk2.Commands
		for name := range Commands {
			commands = append(commands, protocol.Command{
				Name:        "__" + name,
				Description: "GoProxy Command",
			})
		}
		pk2.Commands = commands
	}
	return HandlerContinue
}

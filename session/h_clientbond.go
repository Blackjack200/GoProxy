package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strings"
)

type ClientMovePacket struct {
	AuthInput bool
}

func (c ClientMovePacket) Handle(player *ProxiedPlayer, pk *packet.Packet) bool {
	if c.AuthInput {
		p, ok := (*pk).(*packet.MovePlayer)
		if ok {
			*pk = &packet.PlayerAuthInput{
				Pitch:     p.Pitch,
				Yaw:       p.Yaw,
				Position:  p.Position,
				HeadYaw:   p.Pitch,
				InputData: uint64(player.ClientClientData().CurrentInputMode),
				InputMode: uint32(player.ClientClientData().CurrentInputMode),
				PlayMode:  0,
			}
		}
	}
	return HandlerContinue
}

type PlayerCommandListener struct {
}

func (PlayerCommandListener) Handle(player *ProxiedPlayer, pk *packet.Packet) bool {
	pk2, ok := (*pk).(*packet.CommandRequest)
	if ok {
		//TODO extract command prefix
		if strings.HasPrefix(pk2.CommandLine, "/__") {
			labels := strings.Split(strings.ToLower(strings.TrimPrefix(pk2.CommandLine, "/__")), " ")
			if len(labels) >= 1 {
				for info, command := range Commands {
					if strings.EqualFold(info.Name, labels[0]) {
						go command.Execute(player, labels[1:])
						return HandlerDrop
					}
				}
			}
			return HandlerContinue
		}
	}
	return HandlerContinue
}

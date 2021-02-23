package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ClientMovePacket struct {
	AuthInput bool
}

func (c ClientMovePacket) Handle(s *Session, pk *packet.Packet) bool {
	if c.AuthInput {
		p, ok := (*pk).(*packet.MovePlayer)
		if ok {
			*pk = &packet.PlayerAuthInput{
				Pitch:     p.Pitch,
				Yaw:       p.Yaw,
				Position:  p.Position,
				HeadYaw:   p.Pitch,
				InputData: uint64(s.Client.ClientData().CurrentInputMode),
				InputMode: uint32(s.Client.ClientData().CurrentInputMode),
				PlayMode:  0,
			}
		}
	}
	return HandlerContinue
}

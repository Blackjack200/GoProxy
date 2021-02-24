package session

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type BuiltinClientPacketHandler struct{}

func (BuiltinClientPacketHandler) Handle(s *Session, pk *packet.Packet) bool {
	switch p := (*pk).(type) {
	case *packet.MobEquipment:
		s.Player.HeldSlot = p.HotBarSlot
		s.Player.HeldItem = p.NewItem
	case *packet.MovePlayer:
		s.Player.Position = p.Position
	case *packet.PlayerAuthInput:
		s.Player.Position = p.Position
	}
	return HandlerContinue
}

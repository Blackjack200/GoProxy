package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Handler interface {
	Handle(s *Session, pk *packet.Packet) bool
}

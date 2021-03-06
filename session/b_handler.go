package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Handler interface {
	Handle(player *ProxiedPlayer, pk *packet.Packet) bool
}

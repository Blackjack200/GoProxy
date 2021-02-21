package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ServerBrokenPacketListener struct{}

func (ServerBrokenPacketListener) Handle(_ Session, pk *packet.Packet) bool {
	switch (*pk).(type) {
	case *packet.CraftingData:
		*pk = &packet.CraftingData{}
	case *packet.CreativeContent:
		*pk = &packet.CreativeContent{}
	}
	return false
}

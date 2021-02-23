package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ServerBrokenPacketListener struct{}

func (ServerBrokenPacketListener) Handle(_ *Session, pk *packet.Packet) bool {
	switch p2 := (*pk).(type) {
	case *packet.CraftingData:
		*pk = &packet.CraftingData{}
	case *packet.CreativeContent:
		*pk = &packet.CreativeContent{}
	case *packet.Transfer:
		logrus.Info(p2.Address + ":" + strconv.Itoa(int(p2.Port)))
	}
	return HandlerContinue
}

package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"math"
	"strconv"
)

type SpeedCommand struct {
}

func (at SpeedCommand) Execute(s *Session, args []string) {
	if len(args) >= 1 {
		speed, _ := strconv.ParseFloat(args[0], 32)
		s.Client.WritePacket(&packet.UpdateAttributes{
			EntityRuntimeID: s.Client.GameData().EntityRuntimeID,
			Attributes: []protocol.Attribute{{
				Name:    "minecraft:movement",
				Value:   float32(speed),
				Max:     math.MaxFloat32,
				Min:     0,
				Default: 0.1,
			}},
		})
		s.Client.WritePacket(&packet.UpdateAttributes{
			EntityRuntimeID: s.Client.GameData().EntityRuntimeID,
			Attributes: []protocol.Attribute{{
				Name:    "minecraft:underwater_movement",
				Value:   float32(speed) * 0.2,
				Max:     math.MaxFloat32,
				Min:     0,
				Default: 0.02,
			}},
		})
		s.Client.WritePacket(&packet.UpdateAttributes{
			EntityRuntimeID: s.Client.GameData().EntityRuntimeID,
			Attributes: []protocol.Attribute{{
				Name:    "minecraft:lava_movement",
				Value:   float32(speed) * 0.2,
				Max:     math.MaxFloat32,
				Min:     0,
				Default: 0.02,
			}},
		})
		SendMessage(s.Client, "[Speed] Set to "+args[0])
	}
}
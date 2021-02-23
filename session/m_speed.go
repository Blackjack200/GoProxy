package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"math"
	"strconv"
)

type SpeedCommand struct {
}

func (at SpeedCommand) Execute(s *Session, args []string) bool {
	if len(args) >= 2 {
		speed, _ := strconv.ParseFloat(args[1], 32)
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
		SendMessage(s.Client, "[Speed] Set to "+args[1])
	}
	return true
}

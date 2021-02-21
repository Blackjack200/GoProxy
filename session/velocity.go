package session

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Velocity struct {
	Enable bool
}

func (v Velocity) Handle(s Session, pk *packet.Packet) bool {
	if v.Enable {
		switch pk2 := (*pk).(type) {
		case *packet.SetActorMotion:
			pk2.Velocity = mgl32.Vec3{0, 0, 0}
		}
	}
	return false
}

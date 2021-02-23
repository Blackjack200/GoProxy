package session

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Velocity struct {
	Enable bool
}

func (v Velocity) Handle(_ *Session, pk *packet.Packet) bool {
	if v.Enable {
		switch pk2 := (*pk).(type) {
		case *packet.SetActorMotion:
			pk2.Velocity = mgl32.Vec3{0, 0, 0}
		}
	}
	return HandlerContinue
}

type VelocityCommand struct {
}

func (VelocityCommand) Execute(s *Session, _ []string) bool {
	handler, ok := s.ServerPacketRewriter["velocity"].(*Velocity)
	if ok {
		handler.Enable = !handler.Enable
		f := "Enable"
		if !handler.Enable {
			f = "Disable"
		}
		SendMessage(s.Client, "[Velocity] "+f)
	}
	return true
}

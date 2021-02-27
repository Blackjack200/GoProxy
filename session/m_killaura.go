package session

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"math"
	"strconv"
)

type KillAura struct {
	Enable bool
	Rate   uint8
}

func (k KillAura) Handle(s *Session, pk *packet.Packet) bool {
	p, ok := (*pk).(*packet.MovePlayer)
	if ok && k.Enable &&
		s.Server.GameData().EntityRuntimeID != p.EntityRuntimeID {
		if distance(s.Player.Position, p.Position) <= 12 {
			go func() {
				for i := uint8(0); i < k.Rate; i++ {
					s.Server.WritePacket(s.Attack(p.EntityRuntimeID))
				}
			}()
		}
	}
	return HandlerContinue
}

type KillAuraCommand struct{}

func (KillAuraCommand) Execute(s *Session, args []string) {
	k, ok := s.ServerPacketRewriter["killaura"].(*KillAura)
	if ok {
		if len(args) == 0 {
			k.Enable = !k.Enable
			f := "Enable"
			if !k.Enable {
				f = "Disable"
			}
			SendMessage(s.Client, "[KillAura] "+f)
		}
		if len(args) >= 1 {
			r, _ := strconv.Atoi(args[0])
			k.Rate = uint8(r)
			SendMessage(s.Client, "[KillAura] Rate="+strconv.Itoa(int(k.Rate)))
		}
	}
}

func distance(a mgl32.Vec3, b mgl32.Vec3) float64 {
	return math.Sqrt(
		math.Pow(float64(a.X()-b.X()), 2) +
			math.Pow(float64(a.Y()-b.Y()), 2) +
			math.Pow(float64(a.Z()-b.Z()), 2))
}

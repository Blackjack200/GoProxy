package session

import "github.com/sandertv/gophertunnel/minecraft"

const (
	HandlerDrop     = true
	HandlerContinue = false
)

type Session struct {
	Client               *minecraft.Conn
	Server               *minecraft.Conn
	Translator           *translator
	ClientPacketRewriter map[string]Handler
	ServerPacketRewriter map[string]Handler
}

func initializeClientPacketRewriter(s *Session) {
	s.ClientPacketRewriter = make(map[string]Handler)
	s.ClientPacketRewriter["command"] = &ClientCommandListener{}
	s.ClientPacketRewriter["move"] = &ClientMovePacket{
		s.Server.GameData().ServerAuthoritativeMovementMode != 0,
	}
	s.ClientPacketRewriter["attack"] = &RateAttack{}
	s.ClientPacketRewriter["nofall"] = &NoFall{}
}

func initializeServerPacketRewriter(s *Session) {
	s.ServerPacketRewriter = make(map[string]Handler)
	s.ServerPacketRewriter["broken"] = &ServerBrokenPacketListener{}
	s.ServerPacketRewriter["command"] = &ServerCommandRewrite{}
	s.ServerPacketRewriter["velocity"] = &Velocity{}
}

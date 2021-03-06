package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type NoFall struct {
	Enable bool
}

func (v NoFall) Handle(_ *ProxiedPlayer, pk *packet.Packet) bool {
	if v.Enable {
		switch pk2 := (*pk).(type) {
		case *packet.MovePlayer:
			pk2.OnGround = true
		}
	}
	return HandlerContinue
}

type NoFallCommand struct {
}

func (NoFallCommand) Execute(player *ProxiedPlayer, _ []string) {
	handler, ok := player.Session.ClientPacketRewriter["nofall"].(*NoFall)
	if ok {
		handler.Enable = !handler.Enable
		f := "Enable"
		if !handler.Enable {
			f = "Disable"
		}

		player.sendMessage("[NoFall] " + f)
	}
}

package session

import (
	"errors"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"golang.org/x/oauth2"
	"sync"
)

const (
	HandlerDrop     = true
	HandlerContinue = false
)

var Players sync.Map

func InitializeSessions() {
	Players = sync.Map{}
}

func initializePacketTransfer(player *ProxiedPlayer) {
	go func() {
		s := player.ServerConn()
		defer Close(s)
		//real client <- remote server
		for s == player.ServerConn() {
			pk, err := s.ReadPacket()
			if err != nil {
				return
			}
			player.Session.Translator.translatePacket(&pk)
			if !handlePacket(player.Session.ServerPacketRewriter, player, &pk) {
				if err := player.WritePacketToClient(pk); err != nil {
					if err2, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
						panic(err2)
					}
				}
			}
		}
	}()
}

func initializeClientPacketRewriter(player *ProxiedPlayer) {
	player.Session.ClientPacketRewriter = make(map[string]Handler)
	player.Session.ClientPacketRewriter["command"] = &PlayerCommandListener{}
	player.Session.ClientPacketRewriter["move"] = &ClientMovePacket{
		player.ServerConn().GameData().ServerAuthoritativeMovementMode != 0,
	}
	player.Session.ClientPacketRewriter["attack"] = &RateAttack{}
	player.Session.ClientPacketRewriter["nofall"] = &NoFall{}
}

func initializeServerPacketRewriter(player *ProxiedPlayer) {
	player.Session.ServerPacketRewriter = make(map[string]Handler)
	player.Session.ServerPacketRewriter["server"] = &ServerPacketHandler{}
	player.Session.ServerPacketRewriter["broken"] = &ServerBrokenPacketListener{}
	player.Session.ServerPacketRewriter["command"] = &ServerCommandListener{}
	player.Session.ServerPacketRewriter["velocity"] = &Velocity{}
	player.Session.ServerPacketRewriter["killaura"] = &KillAura{
		Enable: false,
		Rate:   10,
	}
}

func handlePacket(rewriter map[string]Handler, player *ProxiedPlayer, pk *packet.Packet) bool {
	drop := false
	for _, handler := range rewriter {
		drop = handler.Handle(player, pk) || drop
	}
	return drop
}

func NewSession(conn *minecraft.Conn, token oauth2.TokenSource, remote string, bypassResourcePacket bool) *ProxiedPlayer {
	if err := conn.StartGame(conn.GameData()); err != nil {
		panic(err)
	}

	player := newPlayer(&NetworkSession{
		Client: conn,
	}, token, bypassResourcePacket)

	player.Session.Translator = newTranslator(conn.GameData())
	con, _ := Connect(conn, token, remote, bypassResourcePacket)
	if con == nil {
		panic("Failed to Connect to " + remote)
	}
	player.Transfer(con)

	go func() {
		//real client -> remote server
		defer player.close("Quit", false, true)
		for {
			pk, err := player.ClientConn().ReadPacket()
			if err != nil {
				return
			}
			player.Session.Translator.translatePacket(&pk)
			conn := player.Session.Server
			if !handlePacket(player.Session.ClientPacketRewriter, player, &pk) {
				if conn != nil {
					if err := conn.WritePacket(pk); err != nil {
						if err2, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
							panic(err2)
						}
					}
				}
			}
		}
	}()

	initializeClientPacketRewriter(player)
	initializeServerPacketRewriter(player)

	Players.Store(player.UUID, player)
	return player
}

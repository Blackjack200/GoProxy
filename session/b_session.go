package session

import (
	"errors"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
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

func initializePlayer(player *ProxiedPlayer, safe bool) {
	g := sync.WaitGroup{}
	if safe {
		g.Add(2)
	} else {
		g.Add(3)
		go func() {
			if err := player.ClientConn().StartGame(player.ServerConn().GameData()); err != nil {
				panic(err)
			}
			logrus.Info("Downstream Connected: " + player.ClientConn().IdentityData().DisplayName)
			g.Done()
		}()
	}

	go func() {
		player.clearEntities()
		g.Done()
	}()

	go func() {
		if err := player.ServerConn().DoSpawn(); err != nil {
			panic(err)
		}
		logrus.Info("Upstream Connected: " + player.ClientConn().IdentityData().DisplayName)
		g.Done()
	}()
	g.Wait()

	initializePacketTransfer(player)
}

func initializePacketTransfer(player *ProxiedPlayer) {
	go func() {
		defer player.close("ProxyServer Disconnect", false, false)
		//real client -> remote server
		for {
			pk, err := player.ClientConn().ReadPacket()
			if err != nil {
				return
			}
			player.Session.Translator.translatePacket(&pk)
			if !handlePacket(player.Session.ClientPacketRewriter, player, &pk) {
				if err := player.WritePacketToServer(pk); err != nil {
					if err2, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
						panic(err2)
					}
				}
			}
		}
	}()
	go func() {
		defer player.close("ProxyServer Disconnect", false, true)
		//real client <- remote server
		for {
			pk, err := player.ServerConn().ReadPacket()
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
	player.Session.ClientPacketRewriter["client"] = &PlayerClientPacketHandler{}
	player.Session.ClientPacketRewriter["command"] = &PlayerCommandListener{}
	player.Session.ClientPacketRewriter["move"] = &ClientMovePacket{
		player.ServerConn().GameData().ServerAuthoritativeMovementMode != 0,
	}
	player.Session.ClientPacketRewriter["attack"] = &RateAttack{}
	player.Session.ClientPacketRewriter["nofall"] = &NoFall{}
}

func initializeServerPacketRewriter(player *ProxiedPlayer) {
	player.Session.ServerPacketRewriter = make(map[string]Handler)
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

func NewSession(conn *minecraft.Conn, token oauth2.TokenSource, remote string, bypassResourcePacket bool, safe bool) *ProxiedPlayer {
	var src oauth2.TokenSource = nil
	if token != nil {
		src = token
	}

	dialer := minecraft.Dialer{
		TokenSource:  src,
		ClientData:   conn.ClientData(),
		IdentityData: conn.IdentityData(),
		//Compact Nukkit/PMMP
		EnableClientCache: false,
	}

	server, dialErr := dialer.Dial("raknet", remote)

	if dialErr != nil {
		panic(dialErr)
		return nil
	}

	if bypassResourcePacket {
		server.WritePacket(&packet.ResourcePackClientResponse{
			Response:        0,
			PacksToDownload: nil,
		})
	}

	player := newPlayer(&NetworkSession{
		Client: conn,
		Server: server,
	})

	initializeClientPacketRewriter(player)
	initializeServerPacketRewriter(player)

	initializePlayer(player, safe)

	player.Session.Translator = newTranslator(conn.GameData())

	player.Session.Translator.updateTranslatorData(player.ServerGameData())

	_ = player.WritePacketToClient(&packet.SetDifficulty{Difficulty: uint32(player.ServerGameData().Difficulty)})
	_ = player.WritePacketToClient(&packet.GameRulesChanged{GameRules: player.ServerGameData().GameRules})
	_ = player.WritePacketToClient(&packet.SetPlayerGameType{GameType: player.ServerGameData().PlayerGameMode})
	_ = player.WritePacketToClient(&packet.MovePlayer{Position: player.ServerGameData().PlayerPosition})
	_ = player.WritePacketToClient(&packet.SetSpawnPosition{Position: player.ServerGameData().WorldSpawn})
	_ = player.WritePacketToClient(&packet.SetTime{Time: int32(player.ServerGameData().Time)})

	if player.ClientConn().GameData().Dimension != player.ServerGameData().Dimension {
		_ = player.WritePacketToClient(&packet.ChangeDimension{
			Dimension: player.ServerGameData().Dimension,
			Position:  player.ServerGameData().PlayerPosition,
			Respawn:   false,
		})
	}

	Players.Store(player.UUID, player)
	return player
}

package session

import (
	"errors"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/scylladb/go-set/i64set"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"sync"
)

const (
	HandlerDrop     = true
	HandlerContinue = false
)

func initializeSession(s *Session, safe bool) {
	g := sync.WaitGroup{}
	if safe {
		g.Add(2)
	} else {
		g.Add(3)
		go func() {
			if err := s.Client.StartGame(s.Server.GameData()); err != nil {
				panic(err)
			}
			logrus.Info("Downstream Connected: " + s.Client.IdentityData().DisplayName)
			g.Done()
		}()
	}

	go func() {
		s.Player.Entities.Each(func(id int64) bool {
			s.Client.WritePacket(&packet.RemoveActor{EntityUniqueID: id})
			return true
		})
		g.Done()
	}()

	go func() {
		if err := s.Server.DoSpawn(); err != nil {
			panic(err)
		}
		logrus.Info("Upstream Connected: " + s.Client.IdentityData().DisplayName)
		g.Done()
	}()
	g.Wait()

	initializePacketTransfer(s)
}

func initializePacketTransfer(s *Session) {
	go func() {
		defer s.Close("ProxyServer Disconnect", false, false)
		//real client -> remote server
		for {
			pk, err := s.Client.ReadPacket()
			if err != nil {
				return
			}
			s.Translator.translatePacket(&pk)
			if !handlePacket(s.ClientPacketRewriter, s, &pk) {
				if err := s.Server.WritePacket(pk); err != nil {
					if _, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
						return
					}
				}
			}
		}
	}()
	go func() {
		defer s.Close("ProxyServer Disconnect", false, true)
		//real client <- remote server
		for {
			pk, err := s.Server.ReadPacket()
			if err != nil {
				return
			}
			s.Translator.translatePacket(&pk)
			if !handlePacket(s.ServerPacketRewriter, s, &pk) {
				if err := s.Client.WritePacket(pk); err != nil {
					if _, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
						return
					}
				}
			}
		}
	}()
}

func initializeClientPacketRewriter(s *Session) {
	s.ClientPacketRewriter = make(map[string]Handler)
	s.ClientPacketRewriter["client"] = &PlayerClientPacketHandler{}
	s.ClientPacketRewriter["command"] = &PlayerCommandListener{}
	s.ClientPacketRewriter["move"] = &ClientMovePacket{
		s.Server.GameData().ServerAuthoritativeMovementMode != 0,
	}
	s.ClientPacketRewriter["attack"] = &RateAttack{}
	s.ClientPacketRewriter["nofall"] = &NoFall{}
}

func initializeServerPacketRewriter(s *Session) {
	s.ServerPacketRewriter = make(map[string]Handler)
	s.ServerPacketRewriter["broken"] = &ServerBrokenPacketListener{}
	s.ServerPacketRewriter["command"] = &ServerCommandListener{}
	s.ServerPacketRewriter["velocity"] = &Velocity{}
	s.ServerPacketRewriter["killaura"] = &KillAura{
		Enable: false,
		Rate:   10,
	}
}

func handlePacket(rewriter map[string]Handler, s *Session, pk *packet.Packet) bool {
	drop := false
	for _, handler := range rewriter {
		drop = handler.Handle(s, pk) || drop
	}
	return drop
}

func NewSession(conn *minecraft.Conn, token *oauth2.TokenSource, remote string, bypassResourcePacket bool, safe bool) *Session {
	var src oauth2.TokenSource = nil
	if token != nil {
		src = *token
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
		disconnect(conn, "ProxyServer Timeout", true)
		panic(dialErr)
		return nil
	}

	if bypassResourcePacket {
		server.WritePacket(&packet.ResourcePackClientResponse{
			Response:        0,
			PacksToDownload: nil,
		})
	}

	s := &Session{
		Client: conn,
		Server: server,
		Player: ProxiedPlayer{
			HeldSlot: 0,
			Entities: i64set.New(),
		},
	}

	initializeClientPacketRewriter(s)
	initializeServerPacketRewriter(s)

	initializeSession(s, safe)

	s.Translator = newTranslator(conn.GameData())
	s.Translator.updateTranslatorData(server.GameData())

	_ = s.Client.WritePacket(&packet.SetDifficulty{Difficulty: uint32(s.Server.GameData().Difficulty)})
	_ = s.Client.WritePacket(&packet.GameRulesChanged{GameRules: s.Server.GameData().GameRules})
	_ = s.Client.WritePacket(&packet.SetPlayerGameType{GameType: s.Server.GameData().PlayerGameMode})
	_ = s.Client.WritePacket(&packet.MovePlayer{Position: s.Server.GameData().PlayerPosition})
	_ = s.Client.WritePacket(&packet.SetSpawnPosition{Position: s.Server.GameData().WorldSpawn})
	_ = s.Client.WritePacket(&packet.SetTime{Time: int32(s.Server.GameData().Time)})
	if s.Client.GameData().Dimension != s.Server.GameData().Dimension {
		_ = s.Client.WritePacket(&packet.ChangeDimension{
			Dimension: s.Server.GameData().Dimension,
			Position:  s.Server.GameData().PlayerPosition,
			Respawn:   false,
		})
	}

	return s
}

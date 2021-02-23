package session

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"sync"
	"time"
)

func NewSession(conn *minecraft.Conn, token *oauth2.Token, remote string, bypassResourcePacket bool) *Session {
	var ts oauth2.TokenSource = nil
	if token != nil {
		ts = oauth2.StaticTokenSource(token)
	}

	dialer := minecraft.Dialer{
		TokenSource:  ts,
		ClientData:   conn.ClientData(),
		IdentityData: conn.IdentityData(),
		//Compact Nukkit/PMMP
		EnableClientCache: false,
	}

	server, de := dialer.Dial("raknet", remote)

	if server == nil {
		disconnect(conn, "ProxyServer Timeout", true)
		panic(de)
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
	}

	initializeClientPacketRewriter(s)
	initializeServerPacketRewriter(s)

	go initializeSession(s)

	s.Translator = newTranslator(conn.GameData())
	s.Translator.updateTranslatorData(server.GameData())

	_ = s.Client.WritePacket(&packet.SetDifficulty{Difficulty: uint32(s.Server.GameData().Difficulty)})
	_ = s.Client.WritePacket(&packet.GameRulesChanged{GameRules: s.Server.GameData().GameRules})
	_ = s.Client.WritePacket(&packet.SetPlayerGameType{GameType: s.Server.GameData().PlayerGameMode})
	_ = s.Client.WritePacket(&packet.MovePlayer{Position: s.Server.GameData().PlayerPosition})
	_ = s.Client.WritePacket(&packet.SetSpawnPosition{Position: s.Server.GameData().WorldSpawn})

	go func() {
		defer disconnect(s.Client, "ProxyServer Error", false)
		//real client -> remote server
		for {
			pk, err := conn.ReadPacket()
			if pk != nil && err == nil {
				if !handlePacket(s.ClientPacketRewriter, s, &pk) {
					s.Translator.translatePacket(&pk)
					_ = s.Server.WritePacket(pk)
				}
			}
		}
	}()
	go func() {
		defer disconnect(s.Client, "ProxyServer Error", false)
		//real client <- remote server
		for {
			pk, err := s.Server.ReadPacket()
			if pk != nil && err == nil {
				if !handlePacket(s.ServerPacketRewriter, s, &pk) {
					s.Translator.translatePacket(&pk)
					_ = s.Client.WritePacket(pk)
				}
			}
		}
	}()
	return s
}

func handlePacket(rewriter map[string]Handler, s *Session, pk *packet.Packet) bool {
	drop := false
	for _, r := range rewriter {
		drop = r.Handle(s, pk) || drop
	}
	return drop
}

func disconnect(conn *minecraft.Conn, msg string, start bool) {
	if start {
		_ = conn.StartGame(conn.GameData())
	}

	conn.WritePacket(&packet.Disconnect{
		HideDisconnectionScreen: false,
		Message:                 msg,
	})
}

func initializeSession(s *Session) {
	g := sync.WaitGroup{}
	g.Add(2)
	go func() {
		_ = s.Client.StartGame(minecraft.GameData{})
		logrus.Info("Downstream Connected: " + s.Client.IdentityData().DisplayName)
		g.Done()
	}()

	go func() {
		_ = s.Server.DoSpawnTimeout(time.Minute)
		logrus.Info("Upstream Connected: " + s.Client.IdentityData().DisplayName)
		g.Done()
	}()
}

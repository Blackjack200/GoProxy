package session

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"sync"
	"time"
)

type Session struct {
	Client *minecraft.Conn
	Server *minecraft.Conn
}

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

	server, _ := dialer.Dial("raknet", remote)

	s := Session{
		Client: conn,
		Server: server,
	}

	if server == nil {
		_ = conn.StartGame(conn.GameData())
		conn.WritePacket(&packet.Disconnect{
			HideDisconnectionScreen: false,
			Message:                 "ProxyServer Timeout",
		})
		return nil
	}

	if bypassResourcePacket {
		server.WritePacket(&packet.ResourcePackClientResponse{
			Response:        0,
			PacksToDownload: nil,
		})
	}

	initializeSession(s)

	trans := newTranslator(conn.GameData())
	trans.updateTranslatorData(server.GameData())

	logrus.Info("Upstream Connected: " + conn.IdentityData().DisplayName)
	go func() {
		defer conn.WritePacket(&packet.Disconnect{
			HideDisconnectionScreen: false,
			Message:                 "Err",
		})
		rewriter := make(map[uint]Handler)
		rewriter[0] = ClientCommandListener{}
		//real client -> remote server
		for {
			pk, err := conn.ReadPacket()
			if pk != nil && err == nil {
				drop := false
				for _, r := range rewriter {
					drop = r.Handle(s, &pk)
				}

				if !drop {
					//trans.translatePacket(&pk)
					_ = s.Server.WritePacket(pk)
				}
			}
		}
	}()
	go func() {
		defer conn.WritePacket(&packet.Disconnect{
			HideDisconnectionScreen: false,
			Message:                 "Err",
		})
		rewriter := make(map[uint]Handler)
		rewriter[0] = ServerBrokenPacketListener{}
		rewriter[1] = ClientCommandRewrite{}
		//real client <- remote server
		for {
			pk, err := s.Server.ReadPacket()
			if pk != nil && err == nil {
				drop := false
				for _, r := range rewriter {
					drop = r.Handle(s, &pk)
				}

				if !drop {
					//trans.translatePacket(&pk)
					_ = s.Client.WritePacket(pk)
				}
			}
		}
	}()
	return &s
}

func initializeSession(s Session) {
	g := sync.WaitGroup{}
	g.Add(2)
	go func() {
		_ = s.Client.StartGameTimeout(s.Server.GameData(), time.Minute)
		logrus.Info("Downstream Connected: " + s.Client.IdentityData().DisplayName)
		g.Done()
	}()
	go func() {
		_ = s.Server.DoSpawnTimeout(time.Minute)
		g.Done()
	}()
	g.Wait()
}

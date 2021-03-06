package server

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"proxy/config"
	"proxy/session"
	"time"
)

var Running = false

func Start() {
	defer logrus.Info("Shutdown Successfully")
	session.InitializeCommand()

	var provider minecraft.ServerStatusProvider = minecraft.NewStatusProvider(config.Motd())

	if !config.LocalStatus() {
		provider, _ = minecraft.NewForeignStatusProvider(config.RemoteStatus())
	}

	listener, err := minecraft.ListenConfig{
		AuthenticationDisabled: config.ProxySideXBL(),
		MaximumPlayers:         1,
		StatusProvider:         provider,
	}.Listen("raknet", config.Bind())

	if err != nil {
		panic(err)
	}

	session.InitializeSessions()
	logrus.Info("Start GoProxy on " + config.Bind() + " -> " + config.Remote())

	Running = true

	src := config.TokenSrc
	if !config.RemoteXBL() {
		src = nil
	}

	go func() {
		for Running {
			conn, acceptErr := listener.Accept()
			if acceptErr == nil {
				go handleConnection(conn.(*minecraft.Conn), config.Remote(), src, config.SafeConnect())
			}
		}
	}()
	for Running {
		time.Sleep(time.Millisecond * 50)
	}
	_ = listener.Close()
}

func handleConnection(conn *minecraft.Conn, remote string, src oauth2.TokenSource, safe bool) {
	if safe {
		go func() {
			_ = conn.StartGame(conn.GameData())
			logrus.Info("Downstream Connected: " + conn.IdentityData().DisplayName)
		}()
	}
	session.NewSession(conn, src, remote, config.BypassResourcePack(), safe)
}

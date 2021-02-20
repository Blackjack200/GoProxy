package main

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"proxy/config"
	"proxy/session"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	cfgErr := config.Initialize()
	tkErr := config.InitializeToken()

	if cfgErr != nil {
		panic(cfgErr)
	}

	if tkErr != nil {
		panic(tkErr)
	}

	logrus.Info("Start GoProxy on " + config.Bind() + " -> " + config.Remote())

	listener, err := minecraft.ListenConfig{
		AuthenticationDisabled: config.XBL(),
		MaximumPlayers:         20,
		StatusProvider:         minecraft.NewStatusProvider(config.Motd()),
	}.Listen("raknet", config.Bind())

	if err != nil {
		panic(err)
	}

	for {
		conn, err2 := listener.Accept()
		if err2 == nil {
			go handleConnection(conn.(*minecraft.Conn), config.Token, config.Remote())
		}
	}
}

func handleConnection(conn *minecraft.Conn, token *oauth2.Token, remote string) {
	session.NewSession(conn, token, remote, config.BypassResourcePack())
}

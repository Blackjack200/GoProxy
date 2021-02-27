package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"proxy/config"
	"proxy/server"
	"syscall"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	if err := config.Initialize(); err != nil {
		panic(err)
	}

	if err := config.InitializeToken(); err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-sigs
		server.Running = false
	}()

	server.Start()
}

package session

import (
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"runtime"
)

type Command interface {
	Execute(s *Session, args []string)
}

type CommandInfo struct {
	Description string
	Name        string
}

var Commands = make(map[*CommandInfo]*Command)
var i = false

func InitializeCommand() {
	if !i {
		Register("status", "Show Proxy status", StatusCommand{})
		Register("ping", "Show Proxy latency", PingCommand{})
		Register("gc", "Do Garbage Collection with proxy", GCCommand{})
		Register("velocity", "Module Velocity", VelocityCommand{})
		Register("attack", "Module RateAttack", RateCommand{})
		Register("nofall", "Module NoFall", NoFallCommand{})
		Register("speed", "Module Speed", SpeedCommand{})
		Register("fly", "Module Fly", FlyCommand{})
		Register("noclip", "Module NoClip", NoClipCommand{})
		Register("killaura", "Module KillAura", KillAuraCommand{})
		Register("gamemode", "Module GameMode", GameModeCommand{})
		i = true
	}
}

func Register(name string, description string, command Command) {
	info := CommandInfo{
		Name:        name,
		Description: description,
	}
	Commands[&info] = &command
}

func SendMessage(conn *minecraft.Conn, message string) {
	conn.WritePacket(&packet.Text{
		TextType:   packet.TextTypeChat,
		SourceName: "GoProxy",
		Message:    text.Colourf(message),
	})
}

type StatusCommand struct {
}

func (s2 StatusCommand) Execute(s *Session, _ []string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	SendMessage(s.Client, fmt.Sprintf("<green>%d</green>MB", m.Alloc/1024/1024))
}

type GCCommand struct {
}

func (s2 GCCommand) Execute(s *Session, _ []string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.Alloc
	runtime.GC()
	runtime.ReadMemStats(&m)
	SendMessage(s.Client, fmt.Sprintf("Free: <green>%d</green>MB", (before-m.Alloc)/1024/1024))
}

type PingCommand struct{}

func (PingCommand) Execute(s *Session, _ []string) {
	SendMessage(s.Client, fmt.Sprintf("Proxy Ping: <green>%d</green>ms", s.Client.Latency().Milliseconds()))
	SendMessage(s.Client, fmt.Sprintf("Server Ping: <green>%d</green>ms", s.Server.Latency().Milliseconds()))
}

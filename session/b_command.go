package session

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"runtime"
	"strconv"
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
		Message:    message,
	})
}

type StatusCommand struct {
}

func (s2 StatusCommand) Execute(s *Session, _ []string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	SendMessage(s.Client, strconv.FormatUint(m.Alloc/1024/1024, 10)+"MB")
}

type GCCommand struct {
}

func (s2 GCCommand) Execute(s *Session, _ []string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.Alloc
	runtime.GC()
	runtime.ReadMemStats(&m)
	SendMessage(s.Client, "Free: "+strconv.FormatUint((before-m.Alloc)/1024/1024, 10)+"MB")
}

type PingCommand struct{}

func (PingCommand) Execute(s *Session, _ []string) {
	SendMessage(s.Client, "Ping: "+strconv.FormatInt(s.Client.Latency().Milliseconds(), 10)+"ms")
}

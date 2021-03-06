package session

import (
	"fmt"
	"runtime"
)

type Command interface {
	Execute(player *ProxiedPlayer, args []string)
}

type CommandInfo struct {
	Description string
	Name        string
}

var Commands = make(map[CommandInfo]Command)
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
	Commands[info] = command
}

type StatusCommand struct {
}

func (StatusCommand) Execute(player *ProxiedPlayer, _ []string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	player.sendMessage(fmt.Sprintf("<green>%d</green>MB", m.Alloc/1024/1024))
}

type GCCommand struct {
}

func (GCCommand) Execute(player *ProxiedPlayer, _ []string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.Alloc
	runtime.GC()
	runtime.ReadMemStats(&m)
	player.sendMessage(fmt.Sprintf("Free: <green>%d</green>MB", (before-m.Alloc)/1024/1024))
}

type PingCommand struct{}

func (PingCommand) Execute(player *ProxiedPlayer, _ []string) {
	player.sendMessage(fmt.Sprintf("Proxy Ping: <green>%d</green>ms", player.ClientConn().Latency().Milliseconds()))
	player.sendMessage(fmt.Sprintf("Server Ping: <green>%d</green>ms", player.ServerConn().Latency().Milliseconds()))
}

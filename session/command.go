package session

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"runtime"
	"strconv"
)

type Command interface {
	Execute(s *Session, args []string) bool
}

func sendMessage(conn *minecraft.Conn, message string) {
	conn.WritePacket(&packet.Text{
		TextType:   packet.TextTypeChat,
		SourceName: "GoProxy",
		Message:    message,
	})
}

type Status struct {
}

func (s2 Status) Execute(s *Session, args []string) bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	sendMessage(s.Client, strconv.FormatUint(m.Alloc/1024/1024, 10)+"MB")
	return true
}

type Gc struct {
}

func (s2 Gc) Execute(s *Session, args []string) bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.Alloc
	runtime.GC()
	runtime.ReadMemStats(&m)
	sendMessage(s.Client, "Free: "+strconv.FormatUint((before-m.Alloc)/1024/1024, 10)+"MB")
	return true
}

type VelocityCommand struct {
}

func (VelocityCommand) Execute(s *Session, args []string) bool {
	switch handler := s.ServerPacketRewriter["velocity"].(type) {
	case *Velocity:
		handler.Enable = !handler.Enable
		f := "Enable"
		if !handler.Enable {
			f = "Disable"
		}

		sendMessage(s.Client, "[Velocity] "+f)
	}
	return true
}

type AttackCommand struct {
}

func (at AttackCommand) Execute(s *Session, args []string) bool {
	if len(args) >= 2 {
		switch handler := s.ClientPacketRewriter["attack"].(type) {
		case *DoubleAttack:
			r, _ := strconv.Atoi(args[1])
			handler.Repeat = r

			sendMessage(s.Client, "[Attack] Set to "+args[1])
		}
	}
	return true
}

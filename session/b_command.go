package session

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Command interface {
	Execute(s *Session, args []string) bool
}

func SendMessage(conn *minecraft.Conn, message string) {
	conn.WritePacket(&packet.Text{
		TextType:   packet.TextTypeChat,
		SourceName: "GoProxy",
		Message:    message,
	})
}

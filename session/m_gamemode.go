package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strconv"
)

type GameModeCommand struct{}

func (GameModeCommand) Execute(s *Session, args []string) {
	if len(args) >= 1 {
		l, _ := strconv.Atoi(args[0])
		s.Client.WritePacket(&packet.SetPlayerGameType{GameType: int32(l)})
		SendMessage(s.Client, "[GameMode] Send "+strconv.Itoa(l))
	}
}

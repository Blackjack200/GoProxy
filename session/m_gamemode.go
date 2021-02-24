package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strconv"
)

type GameModeCommand struct{}

func (GameModeCommand) Execute(s *Session, args []string) bool {
	if len(args) >= 2 {
		l, _ := strconv.Atoi(args[1])
		s.Client.WritePacket(&packet.SetPlayerGameType{GameType: int32(l)})
		SendMessage(s.Client, "[GameMode] Send "+strconv.Itoa(l))
	}
	return HandlerContinue
}

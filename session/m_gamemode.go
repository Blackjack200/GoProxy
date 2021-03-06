package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strconv"
)

type GameModeCommand struct{}

func (GameModeCommand) Execute(player *ProxiedPlayer, args []string) {
	if len(args) >= 1 {
		l, _ := strconv.Atoi(args[0])
		_ = player.WritePacketToClient(&packet.SetPlayerGameType{GameType: int32(l)})
		player.sendMessage("[GameMode] Send " + strconv.Itoa(l))
	}
}

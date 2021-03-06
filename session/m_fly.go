package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type FlyCommand struct {
}

func (FlyCommand) Execute(player *ProxiedPlayer, _ []string) {
	_ = player.WritePacketToClient(&packet.AdventureSettings{
		Flags:             packet.AdventureFlagAllowFlight,
		PermissionLevel:   packet.PermissionLevelMember,
		PlayerUniqueID:    player.ClientGameData().EntityUniqueID,
		ActionPermissions: uint32(packet.ActionPermissionBuildAndMine | packet.ActionPermissionDoorsAndSwitched | packet.ActionPermissionOpenContainers | packet.ActionPermissionAttackPlayers | packet.ActionPermissionAttackMobs),
	})
	player.sendMessage("[Fly] Sent!!!")
}

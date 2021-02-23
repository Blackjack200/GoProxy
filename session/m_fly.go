package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type FlyCommand struct {
}

func (FlyCommand) Execute(s *Session, _ []string) bool {
	s.Client.WritePacket(&packet.AdventureSettings{
		Flags:             packet.AdventureFlagAllowFlight,
		PermissionLevel:   packet.PermissionLevelMember,
		PlayerUniqueID:    s.Client.GameData().EntityUniqueID,
		ActionPermissions: uint32(packet.ActionPermissionBuildAndMine | packet.ActionPermissionDoorsAndSwitched | packet.ActionPermissionOpenContainers | packet.ActionPermissionAttackPlayers | packet.ActionPermissionAttackMobs),
	})
	SendMessage(s.Client, "[Fly] Sent!!!")
	return true
}

package session

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type NoClipCommand struct {
}

func (at NoClipCommand) Execute(s *Session, _ []string) bool {
	s.Client.WritePacket(&packet.AdventureSettings{
		Flags:             packet.AdventureFlagNoClip,
		PermissionLevel:   packet.PermissionLevelMember,
		PlayerUniqueID:    s.Client.GameData().EntityUniqueID,
		ActionPermissions: uint32(packet.ActionPermissionBuildAndMine | packet.ActionPermissionDoorsAndSwitched | packet.ActionPermissionOpenContainers | packet.ActionPermissionAttackPlayers | packet.ActionPermissionAttackMobs),
	})
	SendMessage(s.Client, "[NoClip] Sent!")
	return true
}

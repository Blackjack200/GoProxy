package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"math"
	"strconv"
)

type SpeedCommand struct {
}

func (SpeedCommand) Execute(player *ProxiedPlayer, args []string) {
	if len(args) >= 1 {
		speed, _ := strconv.ParseFloat(args[0], 32)
		_ = player.WritePacketToClient(&packet.UpdateAttributes{
			EntityRuntimeID: player.ClientGameData().EntityRuntimeID,
			Attributes: []protocol.Attribute{{
				Name:    "minecraft:movement",
				Value:   float32(speed),
				Max:     math.MaxFloat32,
				Min:     0,
				Default: 0.1,
			}},
		})
		_ = player.WritePacketToClient(&packet.UpdateAttributes{
			EntityRuntimeID: player.ClientGameData().EntityRuntimeID,
			Attributes: []protocol.Attribute{{
				Name:    "minecraft:underwater_movement",
				Value:   float32(speed) * 0.2,
				Max:     math.MaxFloat32,
				Min:     0,
				Default: 0.02,
			}},
		})
		_ = player.WritePacketToClient(&packet.UpdateAttributes{
			EntityRuntimeID: player.ClientGameData().EntityRuntimeID,
			Attributes: []protocol.Attribute{{
				Name:    "minecraft:lava_movement",
				Value:   float32(speed) * 0.2,
				Max:     math.MaxFloat32,
				Min:     0,
				Default: 0.02,
			}},
		})
		player.sendMessage("[Speed] Set to " + args[0])
	}
}

type GameModeCommand struct{}

func (GameModeCommand) Execute(player *ProxiedPlayer, args []string) {
	if len(args) >= 1 {
		l, _ := strconv.Atoi(args[0])
		_ = player.WritePacketToClient(&packet.SetPlayerGameType{GameType: int32(l)})
		player.sendMessage("[GameMode] Send " + strconv.Itoa(l))
	}
}

type HighJumpCommand struct {
}

func (HighJumpCommand) Execute(player *ProxiedPlayer, args []string) {
	if len(args) >= 1 {
		level, _ := strconv.Atoi(args[0])

		_ = player.WritePacketToClient(&packet.MobEffect{
			EntityRuntimeID: player.ClientGameData().EntityRuntimeID,
			Operation:       packet.MobEffectRemove,
			EffectType:      packet.EffectJumpBoost,
			Amplifier:       int32(level),
			Duration:        114514,
		})

		_ = player.WritePacketToClient(&packet.MobEffect{
			EntityRuntimeID: player.ClientGameData().EntityRuntimeID,
			Operation:       packet.MobEffectAdd,
			EffectType:      packet.EffectJumpBoost,
			Amplifier:       int32(level),
			Duration:        114514,
		})

		if level <= 0 {
			_ = player.WritePacketToClient(&packet.MobEffect{
				EntityRuntimeID: player.ClientGameData().EntityRuntimeID,
				Operation:       packet.MobEffectRemove,
				EffectType:      packet.EffectJumpBoost,
				Amplifier:       int32(level),
				Duration:        1,
			})
		}
		player.sendMessage("[HighJump] Sent " + args[0] + "!!")
	}
}

type NoClipCommand struct {
}

func (NoClipCommand) Execute(player *ProxiedPlayer, _ []string) {
	_ = player.WritePacketToClient(&packet.AdventureSettings{
		Flags:             packet.AdventureFlagNoClip,
		PermissionLevel:   packet.PermissionLevelMember,
		PlayerUniqueID:    player.ClientGameData().EntityUniqueID,
		ActionPermissions: uint32(packet.ActionPermissionBuildAndMine | packet.ActionPermissionDoorsAndSwitched | packet.ActionPermissionOpenContainers | packet.ActionPermissionAttackPlayers | packet.ActionPermissionAttackMobs),
	})
	player.sendMessage("[NoClip] Sent!")
}

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

package session

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"golang.org/x/oauth2"
)

func Connect(client *minecraft.Conn, token oauth2.TokenSource, remote string, bypassResourcePacket bool) (*minecraft.Conn, error) {
	dialer := minecraft.Dialer{
		TokenSource:  token,
		ClientData:   client.ClientData(),
		IdentityData: client.IdentityData(),
	}

	server, err := dialer.Dial("raknet", remote)

	if err != nil {
		return nil, err
	}

	if er := server.DoSpawn(); er != nil {
		return nil, er
	}

	if bypassResourcePacket {
		server.WritePacket(&packet.ResourcePackClientResponse{
			Response:        0,
			PacksToDownload: nil,
		})
	}

	return server, nil
}

func Close(conn *minecraft.Conn) {
	if err := conn.Close(); err != nil {
		panic(err)
	}
}

func (player *ProxiedPlayer) Transfer(conn *minecraft.Conn) {
	if temp := player.ServerConn(); temp != nil {
		Close(temp)
	}

	player.Session.Server = conn
	player.Session.Translator.updateTranslatorData(player.ServerGameData())

	_ = player.WritePacketToClient(&packet.SetDifficulty{Difficulty: uint32(player.ServerGameData().Difficulty)})
	_ = player.WritePacketToClient(&packet.GameRulesChanged{GameRules: player.ServerGameData().GameRules})
	_ = player.WritePacketToClient(&packet.SetPlayerGameType{GameType: player.ServerGameData().PlayerGameMode})
	_ = player.WritePacketToClient(&packet.MovePlayer{Position: player.ServerGameData().PlayerPosition})
	_ = player.WritePacketToClient(&packet.MoveActorAbsolute{
		EntityRuntimeID: player.ClientGameData().EntityRuntimeID,
		Flags:           packet.MoveFlagTeleport | packet.MoveFlagOnGround,
		Position:        player.ServerGameData().PlayerPosition,
		Rotation:        mgl32.Vec3{},
	})
	_ = player.WritePacketToClient(&packet.SetSpawnPosition{Position: player.ServerGameData().WorldSpawn})
	_ = player.WritePacketToClient(&packet.SetTime{Time: int32(player.ServerGameData().Time)})

	if player.ClientGameData().Dimension != player.ServerGameData().Dimension {
		_ = player.WritePacketToClient(&packet.ChangeDimension{
			Dimension: player.ServerGameData().Dimension,
			Position:  player.ServerGameData().PlayerPosition,
			Respawn:   false,
		})
	}

	player.clearEntities()
	player.clearEffects()
	player.clearPlayerList()
	player.clearBossBars()
	player.clearScoreboard()

	initializePacketTransfer(player)
}

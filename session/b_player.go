package session

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/scylladb/go-set/b16set"
	"github.com/scylladb/go-set/i32set"
	"github.com/scylladb/go-set/i64set"
	"github.com/scylladb/go-set/strset"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type NetworkSession struct {
	Client               *minecraft.Conn
	Server               *minecraft.Conn
	Translator           *translator
	ClientPacketRewriter map[string]Handler
	ServerPacketRewriter map[string]Handler
}

func (player *ProxiedPlayer) ClientConn() *minecraft.Conn {
	return player.Session.Client
}

func (player *ProxiedPlayer) ServerConn() *minecraft.Conn {
	return player.Session.Server
}

func (player *ProxiedPlayer) ClientClientData() login.ClientData {
	return player.ClientConn().ClientData()
}

func (player *ProxiedPlayer) ClientGameData() minecraft.GameData {
	return player.ClientConn().GameData()
}

func (player *ProxiedPlayer) ClientIdentityData() login.IdentityData {
	return player.ClientConn().IdentityData()
}

func (player *ProxiedPlayer) ServerClientData() login.ClientData {
	return player.ServerConn().ClientData()
}

func (player *ProxiedPlayer) ServerGameData() minecraft.GameData {
	return player.ServerConn().GameData()
}

func (player *ProxiedPlayer) ServerIdentityData() login.IdentityData {
	return player.ServerConn().IdentityData()
}

func (player *ProxiedPlayer) WritePacketToClient(p packet.Packet) error {
	return player.ClientConn().WritePacket(p)
}

func (player *ProxiedPlayer) WritePacketToServer(p packet.Packet) error {
	return player.ServerConn().WritePacket(p)
}

type ProxiedPlayer struct {
	Src                  oauth2.TokenSource
	BypassResourcePacket bool
	UUID                 uuid.UUID
	Session              *NetworkSession
	HeldSlot             byte
	HeldItem             protocol.ItemStack
	Position             mgl32.Vec3

	Entities    *i64set.Set
	Scoreboards *strset.Set
	PlayerList  *b16set.Set
	Effects     *i32set.Set
	BossBars    *i64set.Set
}

func (player *ProxiedPlayer) Attack(EntityRuntimeId uint64) {
	_ = player.WritePacketToServer(&packet.InventoryTransaction{
		LegacyRequestID:    0,
		LegacySetItemSlots: []protocol.LegacySetItemSlot{},
		HasNetworkIDs:      false,
		TransactionData: &protocol.UseItemOnEntityTransactionData{
			TargetEntityRuntimeID: EntityRuntimeId,
			ActionType:            protocol.UseItemOnEntityActionAttack,
			HotBarSlot:            int32(player.HeldSlot),
			HeldItem:              player.HeldItem,
			Position:              player.Position,
			ClickedPosition:       mgl32.Vec3{0, 0, 0},
		},
	})
}

func disconnect(conn *minecraft.Conn, msg string, start bool) {
	if start {
		_ = conn.StartGame(conn.GameData())
	}

	conn.WritePacket(&packet.Disconnect{
		HideDisconnectionScreen: false,
		Message:                 msg,
	})

	_ = conn.Close()
}

func (player *ProxiedPlayer) sendMessage(message string) {
	_ = player.WritePacketToClient(&packet.Text{
		TextType:   packet.TextTypeChat,
		SourceName: "GoProxy",
		Message:    text.Colourf(message),
	})
}

func (player *ProxiedPlayer) close(msg string, start bool, log bool) {
	disconnect(player.ClientConn(), msg, start)
	Close(player.ServerConn())
	Close(player.ClientConn())
	if log {
		logrus.Info("Disconnect: " + player.Session.Client.IdentityData().DisplayName)
		Players.Delete(player.UUID)
	}
}

func (player *ProxiedPlayer) clearEntities() {
	player.Entities.Each(func(id int64) bool {
		_ = player.WritePacketToClient(&packet.RemoveActor{EntityUniqueID: id})
		return true
	})
	player.Entities.Clear()
}

func (player *ProxiedPlayer) clearPlayerList() {
	var entries = make([]protocol.PlayerListEntry, player.PlayerList.Size())
	player.PlayerList.Each(func(uid [16]byte) bool {
		entries = append(entries, protocol.PlayerListEntry{UUID: uid})
		return true
	})

	_ = player.WritePacketToClient(&packet.PlayerList{ActionType: packet.PlayerListActionRemove, Entries: entries})

	player.PlayerList.Clear()
}

// clearEffects flushes the effects map and removes all the effects for the client.
func (player *ProxiedPlayer) clearEffects() {
	player.Effects.Each(func(i int32) bool {
		_ = player.WritePacketToClient(&packet.MobEffect{
			EntityRuntimeID: player.ClientGameData().EntityRuntimeID,
			Operation:       packet.MobEffectRemove,
			EffectType:      i,
		})
		return true
	})

	player.Effects.Clear()
}

// clearBossBars clears all of the boss bars currently visible the client.
func (player *ProxiedPlayer) clearBossBars() {
	player.BossBars.Each(func(b int64) bool {
		_ = player.WritePacketToClient(&packet.BossEvent{
			BossEntityUniqueID: b,
			EventType:          packet.BossEventHide,
		})
		return true
	})

	player.BossBars.Clear()
}

// clearScoreboard clears the current scoreboard visible by the client.
func (player *ProxiedPlayer) clearScoreboard() {
	player.Scoreboards.Each(func(sb string) bool {
		_ = player.WritePacketToClient(&packet.RemoveObjective{ObjectiveName: sb})
		return true
	})

	player.Scoreboards.Clear()
}

func newPlayer(s *NetworkSession, src oauth2.TokenSource, bp bool) *ProxiedPlayer {
	return &ProxiedPlayer{
		Src:                  src,
		BypassResourcePacket: bp,
		UUID:                 uuid.New(),
		Session:              s,
		HeldSlot:             0,
		HeldItem:             protocol.ItemStack{},
		Position:             mgl32.Vec3{},
		Entities:             i64set.New(),
		Scoreboards:          strset.New(),
		PlayerList:           b16set.New(),
		Effects:              i32set.New(),
		BossBars:             i64set.New(),
	}
}

type PlayerClientPacketHandler struct{}

func (PlayerClientPacketHandler) Handle(player *ProxiedPlayer, pk *packet.Packet) bool {
	switch p := (*pk).(type) {
	case *packet.MobEquipment:
		player.HeldSlot = p.HotBarSlot
		player.HeldItem = p.NewItem
	case *packet.MovePlayer:
		player.Position = p.Position
	case *packet.PlayerAuthInput:
		player.Position = p.Position
	case *packet.AddActor:
		player.Entities.Add(p.EntityUniqueID)
	case *packet.AddItemActor:
		player.Entities.Add(p.EntityUniqueID)
	case *packet.AddPainting:
		player.Entities.Add(p.EntityUniqueID)
	case *packet.AddPlayer:
		player.Entities.Add(p.EntityUniqueID)
	case *packet.RemoveActor:
		player.Entities.Remove(p.EntityUniqueID)
	case *packet.RemoveObjective:
		player.Scoreboards.Remove(p.ObjectiveName)
	case *packet.SetDisplayObjective:
		player.Scoreboards.Add(p.ObjectiveName)
	case *packet.BossEvent:
		if p.EventType == packet.BossEventShow {
			player.BossBars.Add(p.BossEntityUniqueID)
		} else if p.EventType == packet.BossEventHide {
			player.BossBars.Remove(p.BossEntityUniqueID)
		}
	case *packet.MobEffect:
		if p.Operation == packet.MobEffectAdd {
			player.Effects.Add(p.EffectType)
		} else if p.Operation == packet.MobEffectRemove {
			player.Effects.Remove(p.EffectType)
		}
	case *packet.PlayerList:
		if p.ActionType == packet.PlayerListActionAdd {
			for _, e := range p.Entries {
				player.PlayerList.Add(e.UUID)
			}
		} else {
			for _, e := range p.Entries {
				player.PlayerList.Remove(e.UUID)
			}
		}
	}
	return HandlerContinue
}

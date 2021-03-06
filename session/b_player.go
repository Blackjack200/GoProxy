package session

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/scylladb/go-set/i64set"
	"github.com/sirupsen/logrus"
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

func (player *ProxiedPlayer) ClientGameData() minecraft.GameData {
	return player.ClientConn().GameData()
}

func (player *ProxiedPlayer) ServerClientData() login.ClientData {
	return player.ServerConn().ClientData()
}

func (player *ProxiedPlayer) ClientClientData() login.ClientData {
	return player.ClientConn().ClientData()
}

func (player *ProxiedPlayer) ServerGameData() minecraft.GameData {
	return player.ServerConn().GameData()
}

func (player *ProxiedPlayer) WritePacketToClient(p packet.Packet) error {
	return player.ClientConn().WritePacket(p)
}

func (player *ProxiedPlayer) WritePacketToServer(p packet.Packet) error {
	return player.ServerConn().WritePacket(p)
}

type ProxiedPlayer struct {
	UUID     uuid.UUID
	Session  *NetworkSession
	HeldSlot byte
	HeldItem protocol.ItemStack
	Position mgl32.Vec3

	Entities *i64set.Set
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
	disconnect(player.Session.Client, msg, start)
	_ = player.Session.Server.Close()
	_ = player.Session.Client.Close()
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

func newPlayer(s *NetworkSession) *ProxiedPlayer {
	return &ProxiedPlayer{
		UUID:     uuid.New(),
		Session:  s,
		HeldSlot: 0,
		HeldItem: protocol.ItemStack{},
		Position: mgl32.Vec3{},
		Entities: i64set.New(),
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
	}
	return HandlerContinue
}

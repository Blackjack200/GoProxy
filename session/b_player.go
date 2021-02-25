package session

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/scylladb/go-set/i64set"
)

type ProxiedPlayer struct {
	HeldSlot byte
	HeldItem protocol.ItemStack
	Position mgl32.Vec3

	Entities *i64set.Set
}

type Session struct {
	Client               *minecraft.Conn
	Server               *minecraft.Conn
	Translator           *translator
	ClientPacketRewriter map[string]Handler
	ServerPacketRewriter map[string]Handler
	Player               ProxiedPlayer
}

func (s Session) Attack(EntityRuntimeId uint64) *packet.InventoryTransaction {
	return &packet.InventoryTransaction{
		LegacyRequestID:    0,
		LegacySetItemSlots: []protocol.LegacySetItemSlot{},
		HasNetworkIDs:      false,
		TransactionData: &protocol.UseItemOnEntityTransactionData{
			TargetEntityRuntimeID: EntityRuntimeId,
			ActionType:            protocol.UseItemOnEntityActionAttack,
			HotBarSlot:            int32(s.Player.HeldSlot),
			HeldItem:              s.Player.HeldItem,
			Position:              s.Player.Position,
			ClickedPosition:       mgl32.Vec3{0, 0, 0},
		},
	}
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

func (s Session) Close(msg string, start bool) {
	disconnect(s.Client, msg, start)
	_ = s.Server.Close()
	_ = s.Client.Close()
}

type PlayerClientPacketHandler struct{}

func (PlayerClientPacketHandler) Handle(s *Session, pk *packet.Packet) bool {
	switch p := (*pk).(type) {
	case *packet.MobEquipment:
		s.Player.HeldSlot = p.HotBarSlot
		s.Player.HeldItem = p.NewItem
	case *packet.MovePlayer:
		s.Player.Position = p.Position
	case *packet.PlayerAuthInput:
		s.Player.Position = p.Position
	case *packet.AddActor:
		s.Player.Entities.Add(p.EntityUniqueID)
	case *packet.AddItemActor:
		s.Player.Entities.Add(p.EntityUniqueID)
	case *packet.AddPainting:
		s.Player.Entities.Add(p.EntityUniqueID)
	case *packet.AddPlayer:
		s.Player.Entities.Add(p.EntityUniqueID)
	case *packet.RemoveActor:
		s.Player.Entities.Remove(p.EntityUniqueID)
	}
	return HandlerContinue
}

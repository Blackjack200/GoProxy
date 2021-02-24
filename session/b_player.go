package session

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ProxiedPlayer struct {
	HeldSlot byte
	HeldItem protocol.ItemStack
	Position mgl32.Vec3
}

type Session struct {
	Client               *minecraft.Conn
	Server               *minecraft.Conn
	Translator           *translator
	ClientPacketRewriter map[string]Handler
	ServerPacketRewriter map[string]Handler
	Player               ProxiedPlayer
}

func Attack(s *Session, EntityRuntimeId uint64) *packet.InventoryTransaction {
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

package session

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strconv"
	"time"
)

type KillAura struct {
	Enable bool
	Rate   uint8

	expire int64
}

func (k KillAura) Handle(s *Session, pk *packet.Packet) bool {
	p, ok := (*pk).(*packet.MovePlayer)
	if ok && k.Enable && s.Server.GameData().EntityRuntimeID != p.EntityRuntimeID {
		if k.expire-time.Now().Unix() != -1 {
			s.Server.WritePacket(&packet.InventoryTransaction{
				LegacyRequestID:    0,
				LegacySetItemSlots: []protocol.LegacySetItemSlot{},
				HasNetworkIDs:      false,
				TransactionData: &protocol.UseItemOnEntityTransactionData{
					TargetEntityRuntimeID: p.EntityRuntimeID,
					ActionType:            protocol.UseItemOnEntityActionAttack,
					HotBarSlot:            int32(s.Player.HeldSlot),
					HeldItem:              s.Player.HeldItem,
					Position:              s.Player.Position,
					ClickedPosition:       mgl32.Vec3{0, 0, 0},
				},
			})
		} else {
			k.expire = time.Now().Add(time.Second).Unix()

			s.Server.WritePacket(&packet.PlayerAction{
				EntityRuntimeID: s.Server.GameData().EntityRuntimeID,
				ActionType:      packet.AnimateActionSwingArm,
			})
		}
	}
	return HandlerContinue
}

type KillAuraCommand struct{}

func (KillAuraCommand) Execute(s *Session, args []string) bool {
	k, ok := s.ServerPacketRewriter["killaura"].(*KillAura)
	if ok {
		if len(args) == 1 {
			k.Enable = !k.Enable
			k.expire = time.Now().Add(time.Second).Unix()
			f := "Enable"
			if !k.Enable {
				f = "Disable"
			}
			SendMessage(s.Client, "[KillAura] "+f)
		}
		if len(args) >= 2 {
			r, _ := strconv.Atoi(args[1])
			k.Rate = uint8(r)
			SendMessage(s.Client, "[KillAura] Rate="+strconv.Itoa(int(k.Rate)))
		}
	}
	return HandlerDrop
}

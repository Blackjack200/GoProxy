package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type DoubleAttack struct {
	Repeat int
}

func (t DoubleAttack) Handle(s Session, pk *packet.Packet) bool {
	switch pk2 := (*pk).(type) {
	case *packet.LevelSoundEvent:
		if pk2.SoundType == packet.SoundEventAttackNoDamage {
			return true
		}
	case *packet.InventoryTransaction:
		switch tr := pk2.TransactionData.(type) {
		case *protocol.UseItemOnEntityTransactionData:
			if tr.ActionType == protocol.UseItemOnEntityActionAttack {
				for i := 0; i < t.Repeat; i++ {
					s.Server.WritePacket(pk2)
				}
			}
		}
	}
	return false
}

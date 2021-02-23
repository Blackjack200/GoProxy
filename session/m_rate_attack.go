package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strconv"
)

type RateAttack struct {
	Repeat int
}

func (t RateAttack) Handle(s *Session, pk *packet.Packet) bool {
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
	return HandlerContinue
}

type RateCommand struct {
}

func (at RateCommand) Execute(s *Session, args []string) bool {
	if len(args) >= 2 {
		handler, ok := s.ClientPacketRewriter["attack"].(*RateAttack)
		if ok {
			r, _ := strconv.Atoi(args[1])
			handler.Repeat = r
			SendMessage(s.Client, "[Attack] Set to "+args[1])
		}
	}
	return true
}

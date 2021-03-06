package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strconv"
)

type RateAttack struct {
	Repeat int
}

func (t RateAttack) Handle(player *ProxiedPlayer, pk *packet.Packet) bool {
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
					_ = player.WritePacketToServer(pk2)
				}
			}
		}
	}
	return HandlerContinue
}

type RateCommand struct {
}

func (RateCommand) Execute(player *ProxiedPlayer, args []string) {
	if len(args) >= 1 {
		handler, ok := player.Session.ClientPacketRewriter["attack"].(*RateAttack)
		if ok {
			r, _ := strconv.Atoi(args[0])
			handler.Repeat = r
			player.sendMessage("[Attack] Set to " + args[0])
		}
	}
}

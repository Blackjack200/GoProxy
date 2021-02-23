package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strings"
)

type ClientCommandListener struct {
}

func (ClientCommandListener) Handle(s *Session, pk *packet.Packet) bool {
	pk2, ok := (*pk).(*packet.CommandRequest)
	if ok {
		//TODO extract command prefix
		if strings.HasPrefix(pk2.CommandLine, "/__") {
			labels := strings.Split(strings.ToLower(strings.TrimPrefix(pk2.CommandLine, "/__")), " ")
			if len(labels) >= 1 {
				for name, command := range Commands {
					if strings.EqualFold(name, labels[0]) {
						(*command).Execute(s, labels)
						return HandlerDrop
					}
				}
			}
			return HandlerContinue
		}
	}
	return HandlerContinue
}

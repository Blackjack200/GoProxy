package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strings"
)

type ClientCommandListener struct {
}

func (ClientCommandListener) Handle(s Session, pk *packet.Packet) bool {
	switch pk2 := (*pk).(type) {
	case *packet.CommandRequest:
		//TODO extract command prefix
		if strings.HasPrefix(pk2.CommandLine, "/__") {
			labels := strings.Split(strings.ToLower(strings.TrimPrefix(pk2.CommandLine, "/__")), " ")
			if len(labels) >= 1 {
				for name, command := range Commands {
					if strings.EqualFold(name, labels[0]) {
						command.Execute(&s, labels)
						return true
					}
				}
			}
			return false
		}
	}
	return false
}

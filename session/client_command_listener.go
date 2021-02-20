package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"runtime"
	"strconv"
	"strings"
)

type ClientCommandListener struct{}

func (ClientCommandListener) Handle(s Session, pk *packet.Packet) bool {
	switch pk2 := (*pk).(type) {
	case *packet.CommandRequest:
		if strings.Contains(pk2.CommandLine, "goproxy") {
			s.Client.WritePacket(&packet.Text{
				TextType:         packet.TextTypeChat,
				NeedsTranslation: false,
				SourceName:       "GoProxy",
				Message:          "Hello This is GoProxy",
			})

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			go runtime.GC()
			s.Client.WritePacket(&packet.Text{
				TextType:         packet.TextTypeChat,
				NeedsTranslation: false,
				SourceName:       "GoProxy",
				Message:          "ALLOC_MEM=" + strconv.FormatUint(m.Alloc/1024, 10) + "KB",
			})
			return true
		}
	}
	return false
}

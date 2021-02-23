package session

import (
	"runtime"
	"strconv"
)

type Status struct {
}

func (s2 Status) Execute(s *Session, _ []string) bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	SendMessage(s.Client, strconv.FormatUint(m.Alloc/1024/1024, 10)+"MB")
	return true
}

type Gc struct {
}

func (s2 Gc) Execute(s *Session, _ []string) bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.Alloc
	runtime.GC()
	runtime.ReadMemStats(&m)
	SendMessage(s.Client, "Free: "+strconv.FormatUint((before-m.Alloc)/1024/1024, 10)+"MB")
	return true
}

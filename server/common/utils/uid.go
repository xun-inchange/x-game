package utils

import (
	"sync"
	"time"
)

var Uid *uid

type uid struct {
	sync.Mutex
}

func init() {
	Uid = &uid{}
}

func (m *uid) NextUID() uint64 {
	m.Lock()
	defer m.Unlock()
	return uint64(time.Now().UnixNano())
}

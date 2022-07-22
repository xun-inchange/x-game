package server

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"x-game/server/modles/role"
	network "x-game/x-common/x_net"
)

type severInfo struct {
	users    sync.Map
	ChatMgr  *chatMgr
	NetSever *network.Server
}

var SeverInfo *severInfo

func init() {
	SeverInfo = &severInfo{ChatMgr: initChatMgr()}
}

func (m *severInfo) SetNetServer(s *network.Server) {
	m.NetSever = s
}

func (m *severInfo) StoreUser(user *role.User) {
	m.users.Store(user.Role.AccountId, user)
}

func (m *severInfo) UserLogout(accountId uint64) {
	m.users.Delete(accountId)
	m.ChatMgr.rooms.Range(func(key, value interface{}) bool {
		r := value.(*chatRoom)
		if m.ChatMgr.UserAlreadyJoinRoom(key.(uint64), accountId) {
			r.deleteUser(accountId)
			return false
		}
		return true
	})
}

func (m *severInfo) GetUser(accountId uint64) (*role.User, bool) {
	v, ok := m.users.Load(accountId)
	if !ok {
		return nil, ok
	}
	return v.(*role.User), ok
}

func (m *severInfo) Run() {
	go m.signalListen()
	m.NetSever.Run()
}

func (m *severInfo) stop() {
	m.NetSever.Stop()
	m.ChatMgr.stop()
}

func (m *severInfo) signalListen() {
	log.Printf("start up signal handle!")
	ticker := time.NewTicker(time.Minute * 2)
	sNotify := make(chan os.Signal)
	signal.Notify(sNotify, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case <-sNotify:
			m.stop()
			return
		case <-ticker.C:
			log.Printf("signal handle running!")
		}
	}
}

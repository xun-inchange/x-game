package x_net

import (
	"log"
	"net"
	"sync"
)

type Server struct {
	Addr        string
	Listener    net.Listener
	CSocketConn map[uint64]Conner
}

func NewServer(addr string) *Server {
	m := &Server{}
	m.Addr = addr
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicf("listen addr[%v] err,error is [%v]", addr, err.Error())
	}
	m.Listener = l
	m.CSocketConn = make(map[uint64]Conner)
	return m
}

func (m *Server) Run() {
	for {
		conn, err := m.Listener.Accept()
		if err != nil {
			log.Printf("accept socket err,error is [%v]", err.Error())
			continue
		}
		log.Printf("new a addr[%v] conn", conn.RemoteAddr().String())
		conner := NewConner(conn)
		go conner.Start()
	}
}

func (m *Server) StoreConner(accountId uint64, c Conner) {
	m.CSocketConn[accountId] = c
}

func (m *Server) GetConner(accountId uint64) Conner {
	return m.CSocketConn[accountId]
}

func (m *Server) Stop() {
	err := m.Listener.Close()
	if err != nil {
		log.Printf("server close listen socket err,error is [%v]", err.Error())
	}
	var wg sync.WaitGroup
	for _, c := range m.CSocketConn {
		wg.Add(1)
		go c.Stop()
		wg.Done()
	}
	wg.Wait()
}
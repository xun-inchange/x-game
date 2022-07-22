package c_net

import (
	"log"
	"net"
	"time"
	"x-game/x-common/x_net"
)

func NewClientConner(addr string) x_net.Conner {
	conn := netDial(addr)
	conner := x_net.NewConner(conn)
	return conner
}

func netDial(addr string) net.Conn {
	connectInterval := time.Second
	for {
		conn, err := net.Dial("tcp", addr)
		if conn != nil {
			return conn
		}
		if err != nil {
			log.Printf("dial addr[%v] server err,error is [%v]", addr, err.Error())
		}
		connectInterval *= 3
		time.Sleep(connectInterval)
	}
}

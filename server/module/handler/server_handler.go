package handler

import (
	"log"
	"x-game/proto/server_proto"
	"x-game/x-common/x_net"
)

func ClientHeartBeat(data interface{}, _ x_net.Conner) {
	if _, ok := data.(*server_proto.ClientHeartBeat); !ok {
		return
	}
	log.Printf("receive  socket conn heartbeat")
}

package x_net

import (
	"github.com/golang/protobuf/proto"
	"log"
)

func SendMsg(msgId uint64, msg proto.Message, c Conner) {
	log.Printf("send msgId[%v] msg[%v]", msgId, msg)
	c.(*socketConn).SendMsg(msgId, msg)
}

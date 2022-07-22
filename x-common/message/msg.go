package message

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"log"
	"x-game/x-common/g"
)

type Message struct {
	MsgId uint64 //消息id
	Data  []byte //数据
}

func BytesToMsg(data []byte) *Message {
	return &Message{
		MsgId: decodeMsgId(data[:g.ReadWriteMinLength]),
		Data:  data[g.ReadWriteMinLength:],
	}
}

func MsgToBytes(msgId uint64, msgData proto.Message) []byte {
	dataBytes, err := proto.Marshal(msgData)
	if err != nil {
		log.Printf("proto marshal err,error is [%v]", err.Error())
	}
	idBytes := encodeMsgId(msgId)
	b := bytes.NewBuffer([]byte{})
	b.Write(idBytes)
	b.Write(dataBytes)
	return b.Bytes()
}

func decodeMsgId(msgBytes []byte) uint64 {
	var id uint64
	b := bytes.NewBuffer(msgBytes)
	err := binary.Read(b, binary.BigEndian, &id)
	if err != nil {
		log.Printf("decode msg id error[%v]", err.Error())
	}
	return id
}

func encodeMsgId(msgId uint64) []byte {
	b := bytes.NewBuffer([]byte{})
	err := binary.Write(b, binary.BigEndian, msgId)
	if err != nil {
		log.Printf("encode msg id error[%v]", err.Error())
	}
	return b.Bytes()
}

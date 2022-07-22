package x_net

import (
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"reflect"
	"time"
	"x-game/x-common/g"
	"x-game/x-common/message"
	"x-game/x-common/x_utils"
)

type Conner interface {
	Start()
	Stop()
	ReadMsg()
	WriteMsg()
	SetAccountId(uint64 uint64)
	GetAccountId() uint64
}

type socketConn struct {
	accountId     uint64
	conn          net.Conn
	In            chan *message.Message
	Out           chan []byte
	bySeverCreate bool //是否为服务器监听创建的socket,如果是将设置dead time
	close         chan struct{}
}

func NewConner(c net.Conn) *socketConn {
	m := &socketConn{
		conn:  c,
		In:    make(chan *message.Message, g.MsgLength),
		Out:   make(chan []byte, g.MsgLength),
		close: make(chan struct{}),
	}
	return m
}

func (m *socketConn) SetAccountId(accountId uint64) {
	m.accountId = accountId
}

func (m *socketConn) GetAccountId() uint64 {
	return m.accountId
}

func (m *socketConn) Start() {
	defer x_utils.RecoverErr()
	go m.msgHandle()
	go m.WriteMsg()
	m.ReadMsg()
}

func (m *socketConn) Stop() {
	m.waitMsgHandle()
	m.conn.Close()
	close(m.close)
}

func (m *socketConn) ReadMsg() {
	for {
		buf := make([]byte, g.ReadWriteMaxLength)
		n, err := m.conn.Read(buf)
		if err != nil {
			log.Printf("read msg err,error is [%v]", err.Error())
			m.Stop()
			break
		}
		if n < g.ReadWriteMinLength {
			log.Printf("msg not match condition!")
			continue
		}
		msg := message.BytesToMsg(buf[:n])
		m.In <- msg
		m.setReadDeadLine()
	}
}

func (m *socketConn) WriteMsg() {
	defer x_utils.RecoverErr()
	for {
		select {
		case bytes := <-m.Out:
			_, err := m.conn.Write(bytes)
			if err != nil {
				log.Printf("write msg error[%v]", err.Error())
			}
		case <-m.close:
			break
		}

	}
}

func (m *socketConn) setReadDeadLine() {
	if !m.bySeverCreate {
		return
	}
	m.conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
}

func (m *socketConn) SendMsg(msgId uint64, msg proto.Message) {
	bytes := message.MsgToBytes(msgId, msg)
	m.Out <- bytes
}

func (m *socketConn) msgHandle() {
	defer x_utils.RecoverErr()
	for {
		select {
		case msg, ok := <-m.In:
			if !ok {
				continue
			}
			m.Handle(msg)
		case <-m.close:
			break
		}
	}
}

func (m *socketConn) Handle(msg *message.Message) {
	hm, ok := GetHandlerModel(msg.MsgId)
	if !ok {
		return
	}
	data := reflect.New(hm.T.Elem()).Interface().(proto.Message)
	if err := proto.Unmarshal(msg.Data, data); err != nil {
		log.Printf("Handle err,error is [%v]", err.Error())
		return
	}
	hm.H(data, m)
}

func (m *socketConn) waitMsgHandle() {
	close(m.In)
	for {
		msg, ok := <-m.In
		if !ok {
			break
		}
		m.Handle(msg)
	}
}

package server

import (
	"container/list"
	"github.com/golang/protobuf/proto"
	"log"
	"sync"
	"x-game/proto/chat_proto"
	"x-game/proto/proto_id"
	"x-game/server/common/consts"
	"x-game/x-common/x_net"
	"x-game/x-common/x_utils"
)

type chatMgr struct {
	rooms sync.Map //key-频道类型 value-房间
}

//公共频道房间
type chatRoom struct {
	sync.RWMutex
	channelId      uint64 //频道id
	users          map[uint64]struct{}
	handleMsgCache chan *ChatHandleMsg
	recordsList    *list.List
	close          chan struct{}
}

type ChatHandleMsg struct {
	ChannelId uint64
	AccountId uint64
	Proto     proto.Message
}

func initChatMgr() *chatMgr {
	m := &chatMgr{}
	for _, channelId := range []uint64{consts.AssnChannel, consts.PublicChannel} {
		r := newChatRoom(channelId)
		go r.Start()
		m.storeRoom(r)
	}
	return m
}

func newChatRoom(channelId uint64) *chatRoom {
	m := &chatRoom{}
	m.channelId = channelId
	m.users = make(map[uint64]struct{})
	m.handleMsgCache = make(chan *ChatHandleMsg, 1024)
	m.recordsList = list.New()
	m.close = make(chan struct{})
	return m
}

func (cm *chatMgr) storeRoom(r *chatRoom) {
	cm.rooms.Store(r.channelId, r)
}
func (cm *chatMgr) IsExistRoom(channelId uint64) bool {
	room := cm.getRoom(channelId)
	if room == nil {
		return false
	}
	return true
}
func (cm *chatMgr) getRoom(channelId uint64) *chatRoom {
	v, ok := cm.rooms.Load(channelId)
	if !ok {
		return nil
	}
	return v.(*chatRoom)
}

func (cm *chatMgr) HandleProto(handleMsg *ChatHandleMsg) {
	r := cm.getRoom(handleMsg.ChannelId)
	r.putProtoMsg(handleMsg)
}

func (cm *chatMgr) UserAlreadyJoinRoom(channelId uint64, accountId uint64) bool {
	r := cm.getRoom(channelId)
	return r.userAlreadyJoinRoom(accountId)
}

func (cm *chatMgr) JoinChatRoom(channelId uint64, accountId uint64) {
	r := cm.getRoom(channelId)
	r.appendUser(accountId)
}

func (cm *chatMgr) LeaveChatRoom(channelId uint64, accountId uint64) {
	r := cm.getRoom(channelId)
	r.deleteUser(accountId)
}

func (cm *chatMgr) stop() {
	cm.rooms.Range(func(key, value interface{}) bool {
		r := value.(*chatRoom)
		r.closeHandle()
		return true
	})
}

func (cm *chatMgr) SendRoomsDetail(c x_net.Conner) {
	roomDetailMsg := &chat_proto.ChatRoomsDetailResp{}
	cm.rooms.Range(func(key, value interface{}) bool {
		r := value.(*chatRoom)
		roomDetailMsg.RoomsDetail = append(roomDetailMsg.RoomsDetail, &chat_proto.ChatRoomDetail{ChannelId: key.(uint64), UserNumber: r.usersNumber()})
		return true
	})
	x_net.SendMsg(proto_id.ChatRoomsDetailResp, roomDetailMsg, c)
}

func (m *chatRoom) appendUser(accountId uint64) {
	m.Lock()
	defer m.Unlock()
	m.users[accountId] = struct{}{}
}

func (m *chatRoom) deleteUser(accountId uint64) {
	m.Lock()
	defer m.Unlock()
	delete(m.users, accountId)
}

func (m *chatRoom) userAlreadyJoinRoom(accountId uint64) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.users[accountId]
	return ok
}

func (m *chatRoom) putProtoMsg(handleMsg *ChatHandleMsg) {
	m.handleMsgCache <- handleMsg
}

func (m *chatRoom) Start() {
	defer x_utils.RecoverErr()
	for {
		select {
		case msg := <-m.handleMsgCache:
			m.handleProtoMsg(msg)
		case <-m.close:
			break
		}
	}
}

func (m *chatRoom) handleProtoMsg(handleMsg *ChatHandleMsg) {
	switch handleMsg.Proto.(type) {
	case *chat_proto.ChatReq:
		m.handleUserChat(handleMsg)
	case *chat_proto.ChatRoomChatRecordsReq:
		m.sendChatRecords(handleMsg)
	default:
		log.Printf("handleMsg[%v] don't handle", handleMsg)
	}
}

//处理用户聊天
func (m *chatRoom) handleUserChat(handleMsg *ChatHandleMsg) {
	chatReq := handleMsg.Proto.(*chat_proto.ChatReq)
	m.recordsList.PushBack(chatReq.GetChatMsg())
	c := SeverInfo.NetSever.GetConner(handleMsg.AccountId)
	x_net.SendMsg(proto_id.ChatResp, &chat_proto.ChatResp{SendSuccess: true, ChannelId: m.channelId}, c)
	m.broadMsg(proto_id.ChatNewMsgNotify, &chat_proto.ChatNewMsgNotify{
		ChannelId: m.channelId,
		ChatMsg:   chatReq.ChatMsg,
	})
}

func (m *chatRoom) sendChatRecords(handleMsg *ChatHandleMsg) {
	resp := &chat_proto.ChatRoomChatRecordsResp{
		Result:    consts.ChatSuccess,
		ChannelId: m.channelId,
	}
	for i := m.recordsList.Front(); i != nil; i = i.Next() {
		chatMsg := i.Value.(*chat_proto.ChatMsg)
		resp.MsgRecords = append(resp.MsgRecords, &chat_proto.ChatMsg{SendAccountId: chatMsg.GetSendAccountId(), Content: chatMsg.Content, SendTime: chatMsg.SendTime})
	}
	c := SeverInfo.NetSever.GetConner(handleMsg.AccountId)
	x_net.SendMsg(proto_id.ChatRoomChatRecordsResp, resp, c)
}

//广播消息
func (m *chatRoom) broadMsg(msgID uint64, msg proto.Message) {
	for accountId := range m.users {
		c := SeverInfo.NetSever.GetConner(accountId)
		if c == nil {
			continue
		}
		x_net.SendMsg(msgID, msg, c)
	}
}

func (m *chatRoom) closeHandle() {
	close(m.close)
}

func (m *chatRoom) usersNumber() uint64 {
	m.RLock()
	defer m.RUnlock()
	return uint64(len(m.users))
}

package handler

import (
	"log"
	"x-game/proto/chat_proto"
	"x-game/proto/proto_id"
	"x-game/server/common/consts"
	"x-game/server/modles/server"
	"x-game/x-common/x_net"
)

func ChatReqHandler(data interface{}, c x_net.Conner) {
	req := data.(*chat_proto.ChatReq)
	if _,ok :=server.SeverInfo.GetUser(c.GetAccountId());!ok{
		log.Printf("get acccountId[%v] user is nil",c.GetAccountId())
		return
	}
	if !server.SeverInfo.ChatMgr.IsExistRoom(req.ChannelId) { //不存在房间
		x_net.SendMsg(proto_id.ChatResp, &chat_proto.ChatResp{}, c)
		return
	}
	server.SeverInfo.ChatMgr.HandleProto(&server.ChatHandleMsg{ChannelId: req.GetChannelId(), AccountId: c.GetAccountId(), Proto: req})
}

func ChatRoomsDetailReqHandler(data interface{}, c x_net.Conner) {
	_, ok := data.(*chat_proto.ChatRoomsDetailReq)
	if !ok {
		return
	}
	if _,ok :=server.SeverInfo.GetUser(c.GetAccountId());!ok{
		log.Printf("get acccountId[%v] user is nil",c.GetAccountId())
		return
	}
	server.SeverInfo.ChatMgr.SendRoomsDetail(c)
}

func ChatEnterRoomReqHandler(data interface{}, c x_net.Conner) {
	req := data.(*chat_proto.ChatEnterRoomReq)
	if _,ok :=server.SeverInfo.GetUser(c.GetAccountId());!ok{
		log.Printf("get acccountId[%v] user is nil",c.GetAccountId())
		return
	}
	if !server.SeverInfo.ChatMgr.IsExistRoom(req.GetChannelId()) { //不存在房间
		x_net.SendMsg(proto_id.ChatEnterRoomResp, &chat_proto.ChatEnterRoomResp{ChannelId: req.ChannelId, Result: consts.ChatChannelNotExist}, c)
		return
	}
	if server.SeverInfo.ChatMgr.UserAlreadyJoinRoom(req.GetChannelId(), c.GetAccountId()) { //房间存在用户
		x_net.SendMsg(proto_id.ChatEnterRoomResp, &chat_proto.ChatEnterRoomResp{ChannelId: req.ChannelId, Result: consts.ChatChannelUserAlreadyJoin}, c)
		return
	}
	server.SeverInfo.ChatMgr.JoinChatRoom(req.GetChannelId(), c.GetAccountId())
	x_net.SendMsg(proto_id.ChatEnterRoomResp, &chat_proto.ChatEnterRoomResp{ChannelId: req.ChannelId, Result: consts.ChatSuccess}, c)
}

func ChatLeaveRoomReqHandler(data interface{}, c x_net.Conner) {
	req := data.(*chat_proto.ChatLeaveRoomReq)
	if _,ok :=server.SeverInfo.GetUser(c.GetAccountId());!ok{
		log.Printf("get acccountId[%v] user is nil",c.GetAccountId())
		return
	}
	if !server.SeverInfo.ChatMgr.IsExistRoom(req.GetChannelId()) { //不存在房间
		x_net.SendMsg(proto_id.ChatLeaveRoomResp, &chat_proto.ChatLeaveRoomResp{ChannelId: req.GetChannelId(), Result: consts.ChatChannelNotExist}, c)
		return
	}
	if !server.SeverInfo.ChatMgr.UserAlreadyJoinRoom(req.GetChannelId(), c.GetAccountId()) { //房间不存在用户
		x_net.SendMsg(proto_id.ChatLeaveRoomResp, &chat_proto.ChatLeaveRoomResp{ChannelId: req.GetChannelId(), Result: consts.ChatUserNotJoin}, c)
		return
	}
	server.SeverInfo.ChatMgr.LeaveChatRoom(req.GetChannelId(), c.GetAccountId())
	x_net.SendMsg(proto_id.ChatLeaveRoomResp, &chat_proto.ChatLeaveRoomResp{ChannelId: req.GetChannelId(), Result: consts.ChatSuccess}, c)
}

func ChatRoomChatRecordsReqHandler(data interface{}, c x_net.Conner) {
	req, ok := data.(*chat_proto.ChatRoomChatRecordsReq)
	if !ok {
		return
	}
	if _,ok :=server.SeverInfo.GetUser(c.GetAccountId());!ok{
		log.Printf("get acccountId[%v] user is nil",c.GetAccountId())
		return
	}
	if !server.SeverInfo.ChatMgr.IsExistRoom(req.GetChannelId()) { //不存在房间
		x_net.SendMsg(proto_id.ChatRoomChatRecordsResp, &chat_proto.ChatRoomChatRecordsResp{Result: consts.ChatChannelNotExist, ChannelId: req.GetChannelId()}, c)
		return
	}
	server.SeverInfo.ChatMgr.HandleProto(&server.ChatHandleMsg{ChannelId: req.GetChannelId(), AccountId: c.GetAccountId(), Proto: req})
}

package handler

import (
	"x-game/proto/account_proto"
	"x-game/proto/chat_proto"
	"x-game/proto/proto_id"
	"x-game/x-common/x_net"
)

func InitHandlers() {
	x_net.RegisterHandler(proto_id.AccountLoginReq, &account_proto.AccountLoginReq{}, AccountLoginReqHandler)
	x_net.RegisterHandler(proto_id.AccountLogoutReq, &account_proto.AccountLogoutReq{}, AccountLogoutReqHandler)
	x_net.RegisterHandler(proto_id.ChatReq, &chat_proto.ChatReq{}, ChatReqHandler)
	x_net.RegisterHandler(proto_id.ChatRoomsDetailReq, &chat_proto.ChatRoomsDetailReq{}, ChatRoomsDetailReqHandler)
	x_net.RegisterHandler(proto_id.ChatEnterRoomReq, &chat_proto.ChatEnterRoomReq{}, ChatEnterRoomReqHandler)
	x_net.RegisterHandler(proto_id.ChatLeaveRoomReq, &chat_proto.ChatLeaveRoomReq{}, ChatLeaveRoomReqHandler)
	x_net.RegisterHandler(proto_id.ChatRoomChatRecordsReq, &chat_proto.ChatRoomChatRecordsReq{}, ChatRoomChatRecordsReqHandler)
}

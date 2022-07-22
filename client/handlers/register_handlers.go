package handlers

import (
	"x-game/proto/account_proto"
	"x-game/proto/chat_proto"
	"x-game/proto/proto_id"
	"x-game/x-common/x_net"
)

func InitHandlers() {
	x_net.RegisterHandler(proto_id.AccountLoginResp, &account_proto.AccountLoginResp{}, AccountLoginRespHandler)
	x_net.RegisterHandler(proto_id.AccountLogoutResp, &account_proto.AccountLogoutResp{}, AccountLogoutRespHandler)
	x_net.RegisterHandler(proto_id.ChatResp, &chat_proto.ChatResp{}, ChatRespHandler)
	x_net.RegisterHandler(proto_id.ChatRoomsDetailResp, &chat_proto.ChatRoomsDetailResp{}, ChatRoomDetailRespHandler)
	x_net.RegisterHandler(proto_id.ChatEnterRoomResp, &chat_proto.ChatEnterRoomResp{}, ChatEnterRoomRespHandler)
	x_net.RegisterHandler(proto_id.ChatLeaveRoomResp, &chat_proto.ChatLeaveRoomResp{}, ChatLeaveRoomRespHandler)
	x_net.RegisterHandler(proto_id.ChatNewMsgNotify, &chat_proto.ChatNewMsgNotify{}, ChatNewMsgNotifyHandle)
	x_net.RegisterHandler(proto_id.ChatRoomChatRecordsResp, &chat_proto.ChatRoomChatRecordsResp{}, ChatRoomChatRecordsRespHandler)
}

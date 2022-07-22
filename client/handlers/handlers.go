package handlers

import (
	"log"
	"x-game/proto/account_proto"
	"x-game/proto/chat_proto"
	"x-game/x-common/x_net"
)

func AccountLoginRespHandler(data interface{}, _ x_net.Conner) {
	resp, ok := data.(*account_proto.AccountLoginResp)
	if !ok {
		return
	}
	log.Printf("accountId[%v] login server success", resp.AccountId)
}

func AccountLogoutRespHandler(data interface{}, c x_net.Conner) {
	_, ok := data.(*account_proto.AccountLogoutResp)
	if !ok {
		return
	}
	log.Printf("accountId[%v] logout server success", c.GetAccountId())
}

func ChatRespHandler(data interface{}, c x_net.Conner) {
	resp, ok := data.(*chat_proto.ChatResp)
	if !ok {
		return
	}
	if !resp.SendSuccess {
		return
	}
	log.Printf("accountId[%v] send channelId[%v] chat msg success", c.GetAccountId(), resp.ChannelId)
}

func ChatRoomDetailRespHandler(data interface{}, _ x_net.Conner) {
	resp, ok := data.(*chat_proto.ChatRoomsDetailResp)
	if !ok {
		return
	}
	if len(resp.RoomsDetail) == 0 {
		return
	}
	for _, roomDetail := range resp.RoomsDetail {
		log.Printf("channleId[%v] chat room,online  user number [%v]", roomDetail.ChannelId, roomDetail.UserNumber)
	}
}

func ChatEnterRoomRespHandler(data interface{}, c x_net.Conner) {
	resp, ok := data.(*chat_proto.ChatEnterRoomResp)
	if !ok {
		return
	}
	if resp.Result == 2 {
		log.Printf("ChatEnterRoomRespHandler,channelId[%v] chat room not exist", resp.GetChannelId())
		return
	}
	if resp.Result == 3 {
		log.Printf("accountId[%v] already join channelId[%v] chat room", c.GetAccountId(), resp.GetChannelId())
		return
	}
	if resp.Result == 1 {
		log.Printf("accountId[%v] join channelId[%v] chat room success", c.GetAccountId(), resp.GetChannelId())
	}
}

func ChatLeaveRoomRespHandler(data interface{}, c x_net.Conner) {
	resp, ok := data.(*chat_proto.ChatLeaveRoomResp)
	if !ok {
		return
	}
	if resp.Result == 2 {
		log.Printf("ChatLeaveRoomRespHandler,channelId[%v] chat room not exist", resp.GetChannelId())
		return
	}
	if resp.Result == 4 {
		log.Printf("accountId[%v] user no join channelId[%v] chat room", c.GetAccountId(), resp.GetChannelId())
		return
	}
	if resp.Result == 1 {
		log.Printf("accountId[%v] leave channelId[%v] chat room success", c.GetAccountId(), resp.GetChannelId())
	}
}

func ChatNewMsgNotifyHandle(data interface{}, _ x_net.Conner) {
	resp, ok := data.(*chat_proto.ChatNewMsgNotify)
	if !ok {
		return
	}
	log.Printf("time[%v] accountId[%v] user send channelId[%v] new chat msg,content is [%v]",
		resp.ChatMsg.GetSendTime(), resp.ChatMsg.GetSendAccountId(), resp.GetChannelId(), resp.ChatMsg.GetContent())
}

func ChatRoomChatRecordsRespHandler(data interface{}, _ x_net.Conner) {
	resp, ok := data.(*chat_proto.ChatRoomChatRecordsResp)
	if !ok {
		return
	}
	if resp.GetResult() == 2 {
		log.Printf("ChatRoomChatRecordsRespHandler,channelId[%v] chat room not exist", resp.GetChannelId())
		return
	}
	if len(resp.MsgRecords) == 0 {
		log.Printf("channleId[%v] not exist chat records", resp.GetChannelId())
		return
	}
	log.Printf("channelId[%v] room chat records:", resp.GetChannelId())
	for _, chatRecord := range resp.MsgRecords {
		log.Printf("time[%v] accountId[%v] user send channelId[%v] new chat msg,content is [%v]",
			chatRecord.GetSendTime(), chatRecord.GetSendAccountId(), resp.GetChannelId(), chatRecord.GetContent())
	}
}

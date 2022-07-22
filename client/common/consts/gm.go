package consts

import "fmt"

const (
	GmLogin              = "1" //登录
	GmLogout             = "2" //登出
	GmSendChatMsg        = "3" //发送聊天信息
	GmChatRoomsDetail    = "4" //获取所有聊天室详情
	GmJoinChatRoom       = "5" //加入聊天房间
	GmChatLeaveRoom      = "6" //离开聊天室
	GmGetChatRoomRecords = "7" //获取聊天室聊天记录
	GmGetGmExplain       = "8" //获取gm说明
)

var GmExplain = []string{
	fmt.Sprintf("GM[%v] usage:login server                    exsample:%v 8855", GmLogin, GmLogin),
	fmt.Sprintf("GM[%v] usage:logout server                   exsample:%v", GmLogout, GmLogout),
	fmt.Sprintf("GM[%v] usage:send to channle room msg        exsample:%v [channelId] [chat msg]", GmSendChatMsg, GmSendChatMsg),
	fmt.Sprintf("GM[%v] usage:get all chat room detail data   exsample:%v", GmChatRoomsDetail, GmChatRoomsDetail),
	fmt.Sprintf("GM[%v] usage:join chat room                  exsample:%v [channleId]", GmJoinChatRoom, GmJoinChatRoom),
	fmt.Sprintf("GM[%v] usage:leave chat room                 exsample:%v [channleId]", GmChatLeaveRoom, GmChatLeaveRoom),
	fmt.Sprintf("GM[%v] usage:get chatroom chat records       exsample:%v [channleId]", GmGetChatRoomRecords, GmGetChatRoomRecords),
	fmt.Sprintf("GM[%v] usage:print all gm                    exsample:%v", GmGetGmExplain, GmGetGmExplain),
}

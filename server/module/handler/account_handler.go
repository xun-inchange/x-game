package handler

import (
	"log"
	"x-game/proto/account_proto"
	"x-game/proto/proto_id"
	"x-game/server/modles/role"
	"x-game/server/modles/server"
	"x-game/x-common/x_net"
)

func AccountLoginReqHandler(data interface{}, c x_net.Conner) {
	req := data.(*account_proto.AccountLoginReq)
	user, ok := server.SeverInfo.GetUser(req.AccountId)
	if !ok { //不存在 直接创角
		user = role.NewUser(req.AccountId)
		server.SeverInfo.StoreUser(user)
		log.Printf("create accountId[%v] user success", req.AccountId)
	} else {
		log.Printf("accountId[%v] already exist and user login success", req.AccountId)
	}
	c.SetAccountId(req.AccountId)
	server.SeverInfo.NetSever.StoreConner(user.Role.AccountId, c)
	x_net.SendMsg(proto_id.AccountLoginResp, &account_proto.AccountLoginResp{AccountId: user.Role.AccountId}, c)
}

func AccountLogoutReqHandler(data interface{}, c x_net.Conner) {
	_, ok := data.(*account_proto.AccountLogoutReq)
	if !ok {
		return
	}
	_, ok = server.SeverInfo.GetUser(c.GetAccountId())
	if !ok { //不存在
		return
	}
	log.Printf("accountId[%v] logout server", c.GetAccountId())
	server.SeverInfo.UserLogout(c.GetAccountId())
	x_net.SendMsg(proto_id.AccountLogoutResp, &account_proto.AccountLogoutResp{}, c)
}

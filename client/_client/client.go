package _client

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"x-game/client/c_net"
	cConsts "x-game/client/common/consts"
	"x-game/proto/account_proto"
	"x-game/proto/chat_proto"
	"x-game/proto/proto_id"
	"x-game/proto/server_proto"
	"x-game/x-common/g"
	"x-game/x-common/x_net"
	network "x-game/x-common/x_net"
	"x-game/x-common/x_utils"
)

type clientInfo struct {
	AccountId uint64
	Conner    network.Conner
	Gm        chan string
	Close     chan struct{}
}

var Client *clientInfo

func NewClient(addr string) *clientInfo {
	m := &clientInfo{}
	m.Gm = make(chan string, 10)
	m.Close = make(chan struct{})
	conner := c_net.NewClientConner(addr)
	m.Conner = conner
	m.flagsParse()
	m.Conner.SetAccountId(m.AccountId)
	return m
}

func (m *clientInfo) flagsParse() {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.Uint64Var(&m.AccountId, "accountId", 0, "set _client account id")
	_ = fs.Parse(os.Args[1:])
}

func (m *clientInfo) Run() {
	log.Printf("_client[%v] running", m.AccountId)
	go m.gmStart()
	go m.ScanStdin()
	go m.signalHandle()
	m.Conner.Start()
}

//func (m *clientInfo) sendHeart() {
//	defer time.AfterFunc(g.HeartTime, m.sendClientSocketHeart)
//	time.AfterFunc(g.HeartTime, m.sendClientSocketHeart)
//}

func (m *clientInfo) gmStart() {
	heartTicker := time.NewTicker(g.HeartBeatTime)
	for {
		select {
		case gmStr := <-m.Gm:
			m.handleGmId(gmStr)
		case <-heartTicker.C:
			m.sendClientSocketHeart()
		case <-m.Close:
			break
		}
	}
}

func (m *clientInfo) signalHandle() {
	ticker := time.NewTicker(time.Minute * 2)
	sNotify := make(chan os.Signal)
	signal.Notify(sNotify, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case <-sNotify:
			m.logout()
			return
		case <-ticker.C:
			log.Printf("signal listen running")
		}
	}
}

func (m *clientInfo) Stop() {
	m.Conner.Stop()
	close(m.Close)
}

func (m *clientInfo) handleGmId(str string) {
	gmArr := x_utils.SplitSpace(str)
	if len(gmArr) == 0 {
		return
	}
	switch gmArr[0] {
	case cConsts.GmLogin:
		m.login(gmArr)
	case cConsts.GmLogout:
		m.logout()
	case cConsts.GmSendChatMsg:
		m.sendChatMsg(gmArr)
	case cConsts.GmChatRoomsDetail:
		m.chatRoomsDetail()
	case cConsts.GmJoinChatRoom:
		m.joinChatRoom(gmArr)
	case cConsts.GmChatLeaveRoom:
		m.leaveChatRoom(gmArr)
	case cConsts.GmGetChatRoomRecords:
		m.getChatRecords(gmArr)
	case cConsts.GmGetGmExplain:
		m.printGmExplain()
	default:
		log.Printf("gm[%v] don't handle", gmArr[0])
	}
}

func (m *clientInfo) receiveGmStr(str string) {
	m.Gm <- str
}

func (m *clientInfo) ScanStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	log.Printf("wait stdin!")
	m.printGmExplain()
	for scanner.Scan() {
		m.receiveGmStr(scanner.Text())
	}
}

func (m *clientInfo) printGmExplain() {
	for _, gm := range cConsts.GmExplain {
		fmt.Println(gm)
	}
}

//登录
func (m *clientInfo) login(gmArr []string) {
	if m.AccountId == 0 && len(gmArr) < 2 {
		log.Printf("login err,please write accountId")
		return
	}
	req := &account_proto.AccountLoginReq{}
	if m.AccountId != 0 {
		req.AccountId = m.AccountId
	} else {
		req.AccountId = strToUint64(gmArr[1])
	}
	x_net.SendMsg(proto_id.AccountLoginReq, req, m.Conner)
}

//登出
func (m *clientInfo) logout() {
	req := &account_proto.AccountLogoutReq{}
	x_net.SendMsg(proto_id.AccountLogoutReq, req, m.Conner)
}

//发送聊天信息
func (m *clientInfo) sendChatMsg(gmArr []string) {
	if len(gmArr) < 3 { //不符合要求
		log.Printf("send chatMsg gm err,gmArr[%v]", gmArr)
		return
	}
	req := &chat_proto.ChatReq{
		ChannelId: strToUint64(gmArr[1]),
		ChatMsg:   &chat_proto.ChatMsg{SendAccountId: m.AccountId, Content: strings.Join(gmArr[2:], " "), SendTime: uint64(time.Now().Unix())},
	}
	x_net.SendMsg(proto_id.ChatReq, req, m.Conner)
}

//获取所有聊天室详情
func (m *clientInfo) chatRoomsDetail() {
	req := &chat_proto.ChatRoomChatRecordsReq{}
	x_net.SendMsg(proto_id.ChatRoomsDetailReq, req, m.Conner)
}

//加入聊天房间
func (m *clientInfo) joinChatRoom(gmArr []string) {
	if len(gmArr) < 2 {
		log.Printf("send joinChatRoom gm err,gmArr[%v]", gmArr)
		return
	}
	req := &chat_proto.ChatEnterRoomReq{ChannelId: strToUint64(gmArr[1])}
	x_net.SendMsg(proto_id.ChatEnterRoomReq, req, m.Conner)
}

func (m *clientInfo) leaveChatRoom(gmArr []string) {
	if len(gmArr) < 2 {
		log.Printf("send leaveChatRoom gm err,gmArr[%v]", gmArr)
		return
	}
	req := &chat_proto.ChatLeaveRoomReq{ChannelId: strToUint64(gmArr[1])}
	x_net.SendMsg(proto_id.ChatLeaveRoomReq, req, m.Conner)
}

func (m *clientInfo) getChatRecords(gmArr []string) {
	if len(gmArr) < 2 {
		return
	}
	req := &chat_proto.ChatRoomChatRecordsReq{ChannelId: strToUint64(gmArr[1])}
	x_net.SendMsg(proto_id.ChatRoomChatRecordsReq, req, m.Conner)
}

func (m *clientInfo) sendClientSocketHeart() {
	log.Printf("accountId[%v] send heart msg", m.AccountId)
	x_net.SendMsg(proto_id.ClientSocketHeart, &server_proto.ClientHeartBeat{}, m.Conner)
}

func strToUint64(str string) uint64 {
	v, _ := strconv.ParseUint(str, 10, 64)
	return v
}

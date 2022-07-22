package consts

const (
	_             = iota
	PublicChannel // 1-公共频道
	AssnChannel   // 2-公会频道
)

const (
	_                          = iota
	ChatSuccess                //成功
	ChatChannelNotExist        //频道不存在
	ChatChannelUserAlreadyJoin //用户已经加入频道
	ChatUserNotJoin            //用户未加入
)

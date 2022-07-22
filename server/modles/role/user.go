package role

import (
	"x-game/server/common/consts"
	"x-game/x-common/x_utils"
)

type User struct {
	Role *Role
}

func NewUser(accountId uint64) *User {
	m := &User{}
	m.Role = &Role{
		AccountId: accountId,
		Name:      x_utils.RandString(consts.RandNameLength),
		Sex:       x_utils.RandUint32(consts.SexBoy, consts.SexGirl),
		Power:     10000,
	}
	return m
}

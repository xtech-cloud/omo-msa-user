package cache

import (
	"omo.msa.user/config"
	"omo.msa.user/proxy/nosql"
	"time"
)

type BaseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator string
	Operator string
	CreateTime time.Time
	UpdateTime time.Time
}

type cacheContext struct {
	users  []*UserInfo
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}
	cacheCtx.users = make([]*UserInfo, 0, 1000)

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if nil != err {
		return err
	}

	users,err := nosql.GetAllUsers()
	for _, user := range users {
		t := new(UserInfo)
		t.initInfo(user)
		cacheCtx.users = append(cacheCtx.users, t)
	}
	return nil
}


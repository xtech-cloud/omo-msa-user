package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	accounts []*AccountInfo
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}
	cacheCtx.accounts = make([]*AccountInfo, 0, 100)

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if nil != err {
		return err
	}

	return nil
}

func Context() *cacheContext {
	return cacheCtx
}

func (mine *cacheContext)CreateAccount(phone, psw, operator string) (*AccountInfo,error) {
	db := new(nosql.Account)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAccountNextID()
	db.CreatedTime = time.Now()
	db.Name = phone
	db.Passwords = psw
	db.Creator = operator
	err := nosql.CreateAccount(db)
	if err == nil {
		info := new(AccountInfo)
		info.initInfo(db)
		mine.accounts = append(mine.accounts, info)
		return info, nil
	}
	return nil,err
}

func (mine *cacheContext)GetUser(uid string) *UserInfo {
	for _, account := range mine.accounts {
		info := account.GetUser(uid)
		if info != nil {
			return info
		}
	}
	db,err := nosql.GetUser(uid)
	if err == nil {
		account := mine.getAccount(db.Account)
		if account != nil {
			return account.GetUser(uid)
		}
	}
	return nil
}

func (mine *cacheContext)getAccount(uid string) *AccountInfo {
	for _, account := range mine.accounts {
		if account.UID == uid {
			return account
		}
	}
	db,err := nosql.GetAccount(uid)
	if err == nil {
		info := new(AccountInfo)
		info.initInfo(db)
		mine.accounts = append(mine.accounts, info)
		return info
	}
	return nil
}

func (mine *cacheContext)getAccountByName(name string) *AccountInfo {
	for _, account := range mine.accounts {
		if account.Name == name {
			return account
		}
	}
	db,err := nosql.GetAccountByName(name)
	if err == nil {
		info := new(AccountInfo)
		info.initInfo(db)
		mine.accounts = append(mine.accounts, info)
		return info
	}
	return nil
}

func (mine *cacheContext)SignIn(name, psw string) (bool,error) {
	account := mine.getAccountByName(name)
	if account == nil {
		return false, errors.New("not found the account")
	}
	if account.Passwords != psw {
		return false, errors.New("the passwords valid failed")
	}
	return true,nil
}


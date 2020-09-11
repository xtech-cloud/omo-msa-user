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
	Creator    string
	Operator   string
	CreateTime time.Time
	UpdateTime time.Time
}

type cacheContext struct {
	accounts []*AccountInfo
	wechats []*WechatInfo
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

func (mine *cacheContext) CreateAccount(name, psw, operator string) (*AccountInfo, error) {
	if len(name) < 1 {
		return nil, errors.New("the account name is empty")
	}
	account := mine.getAccountByName(name)
	if account != nil {
		return account,nil
	}

	db := new(nosql.Account)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAccountNextID()
	db.CreatedTime = time.Now()
	db.Name = name
	db.Passwords = psw
	db.Creator = operator
	err := nosql.CreateAccount(db)
	if err == nil {
		info := new(AccountInfo)
		info.initInfo(db)
		mine.accounts = append(mine.accounts, info)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetUser(uid string) *UserInfo {
	for _, account := range mine.accounts {
		info := account.GetUser(uid)
		if info != nil {
			return info
		}
	}

	db, err := nosql.GetUser(uid)
	if err == nil {
		account := mine.GetAccount(db.Account)
		if account != nil {
			user := new(UserInfo)
			user.initInfo(db)
			account.Users = append(account.Users, user)
			return user
		}
	}
	return nil
}

func (mine *cacheContext) GetUserByEntity(entity string) *UserInfo {
	for _, account := range mine.accounts {
		info := account.GetUser(entity)
		if info != nil {
			return info
		}
	}
	db, err := nosql.GetUserByEntity(entity)
	if err == nil {
		account := mine.GetAccount(db.Account)
		if account != nil {
			db,err := nosql.GetUserByEntity(entity)
			if err == nil {
				user := new(UserInfo)
				user.initInfo(db)
				account.Users = append(account.Users, user)
				return user
			}
		}
	}
	return nil
}

func (mine *cacheContext) GetUserByPhone(phone string) *UserInfo {
	for _, account := range mine.accounts {
		info := account.GetUserByPhone(phone)
		if info != nil {
			return info
		}
	}
	db, err := nosql.GetUserByPhone(phone)
	if err == nil {
		info := new(UserInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetUserBySNS(uid string, kind uint8) *UserInfo {
	for _, account := range mine.accounts {
		if account.DefaultUser().HadSNS(uid) {
			return account.DefaultUser()
		}
	}
	db, err := nosql.GetUserBySNS(uid)
	if err == nil {
		info := new(UserInfo)
		info.initInfo(db)
		return info
	}
	return nil
}


func (mine *cacheContext) GetAccount(uid string) *AccountInfo {
	for _, account := range mine.accounts {
		if account.UID == uid {
			return account
		}
	}
	db, err := nosql.GetAccount(uid)
	if err == nil {
		info := new(AccountInfo)
		info.initInfo(db)
		mine.accounts = append(mine.accounts, info)
		return info
	}
	return nil
}

func (mine *cacheContext) GetAccountByUser(user string) *AccountInfo {
	for _, account := range mine.accounts {
		if account.HadUser(user) {
			return account
		}
	}
	db, err := nosql.GetUser(user)
	if err == nil {
		account := mine.GetAccount(db.Account)
		if account != nil {
			return account
		}
	}
	return nil
}

func (mine *cacheContext) getAccountByName(name string) *AccountInfo {
	for _, account := range mine.accounts {
		if account.Name == name {
			return account
		}
	}
	db, err := nosql.GetAccountByName(name)
	if err == nil {
		info := new(AccountInfo)
		info.initInfo(db)
		mine.accounts = append(mine.accounts, info)
		return info
	}
	return nil
}

func (mine *cacheContext) SignIn(name, psw string) (string, error) {
	account := mine.getAccountByName(name)
	if account == nil {
		return "", errors.New("not found the account")
	}
	if account.Passwords != psw {
		return "", errors.New("the passwords valid failed")
	}
	return account.DefaultUser().UID, nil
}

func (mine *cacheContext) AllUsers() []*UserInfo {
	list := make([]*UserInfo, 0, 10)
	for _, account := range mine.accounts {
		list = append(list, account.Users...)
	}
	return list
}

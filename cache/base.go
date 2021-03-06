package cache

import (
	"errors"
	"github.com/micro/go-micro/v2/logger"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
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
	DeleteTime time.Time
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
	userNum := nosql.GetUserCount()
	accNum := nosql.GetAccountCount()
	wxNum := nosql.GetWechatCount()
	logger.Infof("the user count = %d; the account count = %d and wechat count = %d", userNum, accNum, wxNum)
	if accNum < 1 {
		account,er := cacheCtx.CreateAccount(config.Schema.Root.Name, config.Schema.Root.Passwords, "system")
		if er != nil {
			return er
		}
		_, er = account.CreateUser(&pb.ReqUserAdd{Name: config.Schema.Root.Name, Operator: "system"})
		if er != nil {
			return er
		}
	}
	return nil
}

func Context() *cacheContext {
	return cacheCtx
}

func (mine *cacheContext) CreateAccount(name, psw, operator string) (*AccountInfo, error) {
	if len(name) < 2 {
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
			user.initInfo(db, account.Status)
			account.Users = append(account.Users, user)
			return user
		}
	}
	return nil
}

func (mine *cacheContext) GetUserByID(id uint64) *UserInfo {
	for _, account := range mine.accounts {
		info := account.GetUserByID(id)
		if info != nil {
			return info
		}
	}

	db, err := nosql.GetUserByID(id)
	if err == nil {
		account := mine.GetAccount(db.Account)
		if account != nil {
			user := new(UserInfo)
			user.initInfo(db, account.Status)
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
				user.initInfo(db, account.Status)
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
		account := mine.GetAccount(db.Account)
		if account != nil {
			info := new(UserInfo)
			info.initInfo(db, account.Status)
			account.Users = append(account.Users, info)
			return info
		}
	}
	return nil
}

func (mine *cacheContext) GetUserBySNS(uid string, kind uint8) *UserInfo {
	for _, account := range mine.accounts {
		if account.DefaultUser() != nil && account.DefaultUser().HadSNS(uid) {
			return account.DefaultUser()
		}
	}
	db, err := nosql.GetUserBySNS(uid)
	if err == nil {
		account := mine.GetAccount(db.Account)
		if account != nil {
			info := new(UserInfo)
			info.initInfo(db, account.Status)
			account.Users = append(account.Users, info)
			return info
		}
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

func (mine *cacheContext)RemoveUser(user, operator string) error {
	for _, account := range mine.accounts {
		for i := 0;i < len(account.Users);i += 1 {
			if account.Users[i].UID == user {
				return account.RemoveUser(user, operator)
			}
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

func removeAccount(uid string) {
	for i := 0;i < len(cacheCtx.accounts);i += 1 {
		if cacheCtx.accounts[i].UID == uid {
			cacheCtx.accounts = append(cacheCtx.accounts[:i], cacheCtx.accounts[i+1:]...)
			break
		}
	}
}

func (mine *cacheContext) SignIn(name, psw string) (string, error) {
	account := mine.getAccountByName(name)
	if account == nil {
		return "", errors.New("not found the account")
	}
	if account.Status == AccountStatusFreeze {
		return "", errors.New("the account is freeze")
	}
	if account.Passwords != psw {
		return "", errors.New("the passwords valid failed")
	}
	account.UpdateTime = time.Now()
	return account.DefaultUser().UID, nil
}

func (mine *cacheContext) AllUsers() []*UserInfo {
	list := make([]*UserInfo, 0, 100)
	for _, account := range mine.accounts {
		list = append(list, account.Users...)
	}
	return list
}

func (mine *cacheContext) SearchUsers(kind pb.UserType, tags []string) []*UserInfo {
	list := make([]*UserInfo, 0, 10)
	ty := uint8(kind)
	if kind < 1 {
		return list
	}
	users,err := nosql.GetUsersByType(ty)
	if err != nil {
		return list
	}
	for _, user := range users {
		if hadKey(user.Tags, tags){
			info := new(UserInfo)
			info.initInfo(user, 0)
			list = append(list, info)
		}
	}
	return list
}

func hadKey(source []string, dest []string) bool {
	if dest == nil || len(dest) < 1 {
		return true
	}
	if source == nil || len(source) < 1 {
		return true
	}
	for _, k := range dest {
		for _, c := range source {
			if k == c {
				return true
			}
		}
	}
	return false
}
package cache

import (
	"errors"
	"fmt"
	"github.com/micro/go-micro/v2/logger"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if nil != err {
		return err
	}
	userNum := nosql.GetUserCount()
	accNum := nosql.GetAccountCount()
	wxNum := nosql.GetWechatCount()
	logger.Infof("the user count = %d; the account count = %d and wechat count = %d", userNum, accNum, wxNum)
	if accNum < 1 {
		psw := CryptPsw(config.Schema.Root.Passwords)
		account,er := cacheCtx.CreateAccount(config.Schema.Root.Name, psw, "system")
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

func fixAllPsw(){
	accounts,_ := nosql.GetAllAccounts()
	for _, account := range accounts {
		if len(account.Passwords) == 32 {
			nosql.UpdateAccountPasswords(account.UID.Hex(), CryptPsw(account.Passwords), "system")
		}
	}
}

func testSignIn()  {
	_,err := cacheCtx.SignIn("18990727722", "25d55ad283aa400af464c76d713c07ad")
	if err != nil {
		fmt.Println(err.Error())
	}else{
		fmt.Println("success")
	}
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
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetUser(uid string) *UserInfo {
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

func CryptPsw(psw string) string {
	if psw == "" {
		return psw
	}
	hash,err := bcrypt.GenerateFromPassword([]byte(psw), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("crypt psw error: " + err.Error())
		return psw
	}
	logger.Info("crypt psw = " + string(hash))
	return string(hash)
}

func (mine *cacheContext) GetUserByID(id uint64) *UserInfo {
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
	db, err := nosql.GetAccount(uid)
	if err == nil {
		info := new(AccountInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetAccountByName(name string) *AccountInfo {
	db, err := nosql.GetAccountByName(name)
	if err == nil {
		info := new(AccountInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetAccountByUser(user string) *AccountInfo {
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
	account := mine.GetAccountByUser(user)
	if account != nil {
		return account.DeleteUser(user)
	}
	return nil
}

func (mine *cacheContext) getAccountByName(name string) *AccountInfo {
	db, err := nosql.GetAccountByName(name)
	if err == nil {
		info := new(AccountInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) SignIn(name, psw string) (string, error) {
	account := mine.getAccountByName(name)
	if account == nil {
		return "", errors.New("not found the account")
	}
	if account.Status == AccountStatusFreeze {
		return "", errors.New("the account is freeze")
	}
	//hash := CryptPsw(psw)
	//hash := "$2a$10$DlGnCj1SSfDaZbxUrXqJ6eEHrKFdgFa8n7MAJODE7zzvrr1SHa6e2"
	err := bcrypt.CompareHashAndPassword([]byte(account.Passwords), []byte(psw))
	if err != nil {
		return "", errors.New("the passwords valid failed that hash = "+account.Passwords+" ; err = " + err.Error())
	}
	account.UpdateTime = time.Now()
	if account.DefaultUser() == nil {
		return "", errors.New("the account not found the user that name = "+name)
	}
	return account.DefaultUser().UID, nil
}

func (mine *cacheContext) AllUsers() []*UserInfo {
	list := make([]*UserInfo, 0, 100)
	all,err := nosql.GetAllUsers()
	if err == nil {
		for _, db := range all {
			acc := mine.GetAccount(db.Account)
			if acc != nil {
				info := new(UserInfo)
				info.initInfo(db, acc.Status)
				list = append(list, info)
			}
		}
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
			account := mine.GetAccount(user.Account)
			if account != nil {
				info := new(UserInfo)
				info.initInfo(user, account.Status)
				list = append(list, info)
			}
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
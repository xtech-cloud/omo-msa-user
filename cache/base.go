package cache

import (
	"errors"
	"fmt"
	"github.com/micro/go-micro/v2/logger"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
	"golang.org/x/crypto/bcrypt"
	"math"
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
		account, er := cacheCtx.CreateAccount(config.Schema.Root.Name, psw, "system")
		if er != nil {
			return er
		}
		_, er = account.CreateUser(&pb.ReqUserAdd{Name: config.Schema.Root.Name, Operator: "system"})
		if er != nil {
			return er
		}
	}
	UpdateMessageType()
	return nil
}

func fixAllPsw() {
	accounts, _ := nosql.GetAllAccounts()
	for _, account := range accounts {
		if len(account.Passwords) == 32 {
			nosql.UpdateAccountPasswords(account.UID.Hex(), CryptPsw(account.Passwords), "system")
		}
	}
}

func testSignIn() {
	_, err := cacheCtx.SignIn("18990727722", "25d55ad283aa400af464c76d713c07ad")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("success")
	}
}

func Context() *cacheContext {
	return cacheCtx
}

func getPageStart(page, num uint32) (int64, int64) {
	var start uint32
	if page < 1 {
		page = 0
		num = 0
		start = 0
	} else {
		if num < 1 {
			num = 10
		}
		start = (page - 1) * num
	}
	return int64(start), int64(num)
}

func CheckPage[T any](page, number uint32, all []T) (uint32, uint32, []T) {
	if len(all) < 1 {
		return 0, 0, make([]T, 0, 1)
	}
	if number < 1 {
		number = 10
	}
	total := uint32(len(all))
	if len(all) <= int(number) {
		return total, 1, all
	}
	maxPage := total/number + 1
	if page < 1 {
		return total, maxPage, all
	}

	var start = (page - 1) * number
	var end = start + number
	if end > total-1 {
		end = total - 1
	}
	list := make([]T, 0, number)
	list = append(all[start:end])
	return total, maxPage, list
}

func CryptPsw(psw string) string {
	if psw == "" {
		return psw
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(psw), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("crypt psw error: " + err.Error())
		return psw
	}
	logger.Info("crypt psw = " + string(hash))
	return string(hash)
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
		return "", errors.New("the passwords valid failed that hash = " + account.Passwords + " ; err = " + err.Error())
	}
	account.UpdateTime = time.Now()
	if account.DefaultUser() == nil {
		return "", errors.New("the account not found the user that name = " + name)
	}
	return account.DefaultUser().UID, nil
}

func (mine *cacheContext) AllUsers() []*UserInfo {
	list := make([]*UserInfo, 0, 100)
	all, err := nosql.GetAllUsers()
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

func (mine *cacheContext) GetUsersByPage(page, num uint32) (uint32, uint32, []*UserInfo) {
	start, number := getPageStart(page, num)
	array, err := nosql.GetUsersByPage(start, number)
	total := nosql.GetUserCount()
	pages := math.Ceil(float64(total) / float64(number))
	if err == nil {
		list := make([]*UserInfo, 0, len(array))
		for _, item := range array {
			info := new(UserInfo)
			info.initInfo(item, 0)
			list = append(list, info)
		}
		return uint32(total), uint32(pages), list
	}
	return 0, 0, make([]*UserInfo, 0, 1)
}

func (mine *cacheContext) GetUsersByPageScene(scene string, page, num uint32) (uint32, uint32, []*UserInfo) {
	start, number := getPageStart(page, num)
	array, err := nosql.GetUsersByScenePage(scene, start, number)
	total := nosql.GetUsersCountByScene(scene)
	pages := math.Ceil(float64(total) / float64(number))
	if err == nil {
		list := make([]*UserInfo, 0, len(array))
		for _, item := range array {
			info := new(UserInfo)
			info.initInfo(item, 0)
			list = append(list, info)
		}
		return uint32(total), uint32(pages), list
	}
	return 0, 0, make([]*UserInfo, 0, 1)
}

func (mine *cacheContext) GetUsersByLatest(scene string, page, num uint32) (uint32, uint32, []*UserInfo) {
	start, number := getPageStart(page, num)
	total := nosql.GetBehavioursCountBy(scene, TargetTypeScene, BehaviourActionRelate)
	pages := math.Ceil(float64(total) / float64(number))
	arr, err := nosql.GetBehavioursByPage(scene, TargetTypeScene, BehaviourActionRelate, start, number)
	if err != nil {
		return 0, 0, nil
	}
	list := make([]*UserInfo, 0, len(arr))
	for _, db := range arr {
		user := mine.GetUser(db.Creator)
		if user != nil {
			list = append(list, user)
		}
	}
	return uint32(total), uint32(pages), list
}

func (mine *cacheContext) SearchUsers(kind pb.UserType, tags []string) []*UserInfo {
	list := make([]*UserInfo, 0, 10)
	ty := uint8(kind)
	if kind < 1 {
		return list
	}
	users, err := nosql.GetUsersByType(ty)
	if err != nil {
		return list
	}
	for _, user := range users {
		if hadKey(user.Tags, tags) {
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

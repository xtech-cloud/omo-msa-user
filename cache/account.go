package cache

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.user/proxy/nosql"
	"time"
)

type AccountInfo struct {
	BaseInfo
	Passwords string
	Users []*UserInfo
}

func (mine *AccountInfo)initInfo(db *nosql.Account)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Name = db.Name
	mine.Passwords = db.Passwords
	users,err := nosql.GetUsersByAccount(mine.UID)
	if err == nil {
		mine.Users = make([]*UserInfo, 0, len(users))
		for _, user := range users {
			info := new(UserInfo)
			info.initInfo(user)
			mine.Users = append(mine.Users, info)
		}
	}else{
		mine.Users = make([]*UserInfo, 0, 1)
	}
}

func (mine *AccountInfo)UpdateName(name, operator string) error {
	err := nosql.UpdateAccountBase(mine.UID, name, operator)
	if err == nil {
		mine.Name = name
		mine.Operator = operator
	}
	return err
}

func (mine *AccountInfo)UpdatePasswords(psw, operator string) error {
	err := nosql.UpdateAccountPasswords(mine.UID, psw, operator)
	if err == nil {
		mine.Passwords = psw
		mine.Operator = operator
	}
	return err
}

func (mine *AccountInfo)CreateUser(name, remark, nick, phone string, tp uint8, sex uint8) (*UserInfo, error) {
	db := new(nosql.User)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetUserNextID()
	db.CreatedTime = time.Now()
	db.Name = name
	db.Type = tp
	db.Account = mine.UID
	db.Remark = remark
	db.Nick = nick
	db.Phone = phone
	db.Sex = sex
	err := nosql.CreateUser(db)
	if err == nil {
		user :=new(UserInfo)
		user.initInfo(db)
		mine.Users = append(mine.Users, user)
		return user,nil
	}
	return nil,err
}

func (mine *AccountInfo)createDatum(info *DatumInfo) error {
	db := new(nosql.Datum)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetDatumNextID()
	db.CreatedTime = time.Now()
	db.Job = info.Job
	err := nosql.CreateDatum(db)
	if err == nil {
		info.initInfo(db)
	}
	return err
}

func (mine *AccountInfo)AllUsers() []*UserInfo {
	return mine.Users
}

func (mine *AccountInfo)GetUser(uid string) *UserInfo {
	for i := 0;i < len(mine.Users);i += 1 {
		if mine.Users[i].UID == uid {
			return mine.Users[i]
		}
	}
	db,err := nosql.GetUser(uid)
	if err == nil {
		user := new(UserInfo)
		user.initInfo(db)
		mine.Users = append(mine.Users, user)
		return user
	}
	return nil
}

func (mine *AccountInfo)GetUserByPhone(phone string) *UserInfo {
	for i := 0;i < len(mine.Users);i += 1 {
		if mine.Users[i].Datum.Phone == phone {
			return mine.Users[i]
		}
	}
	db,err := nosql.GetUserByPhone(phone)
	if err == nil {
		user := new(UserInfo)
		user.initInfo(db)
		mine.Users = append(mine.Users, user)
		return user
	}
	return nil
}

func (mine *AccountInfo)RemoveUser(uid, operator string) error {
	if len(uid) < 1{
		return errors.New("the user uid is empty")
	}
	err := nosql.RemoveUser(uid, operator)
	if err == nil {
		for i := 0;i < len(mine.Users);i += 1 {
			if mine.Users[i].UID == uid {
				mine.Users = append(mine.Users[:i], mine.Users[i+1:]...)
				break
			}
		}
	}
	return err
}
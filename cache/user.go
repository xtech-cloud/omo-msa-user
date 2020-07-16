package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.user/proxy/nosql"
	"time"
)

type UserInfo struct {
	Type uint8
	BaseInfo
	Account string
	Remark string
	Datum DatumInfo
}

func CreateUser(account , name, remark string, tp uint8, info *DatumInfo) (*UserInfo, error) {
	err1 := createDatum(info)
	if err1 != nil {
		return nil, err1
	}
	db := new(nosql.User)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetUserNextID()
	db.CreatedTime = time.Now()
	db.Name = name
	db.Type = tp
	db.Account = account
	db.Remark = remark
	db.Datum = info.UID
	err := nosql.CreateUser(db)
	if err == nil {
		user :=new(UserInfo)
		user.initInfo(db)
		cacheCtx.users = append(cacheCtx.users, user)
		return user,nil
	}
	return nil,err
}

func createDatum(info *DatumInfo) error {
	db := new(nosql.Datum)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetDatumNextID()
	db.CreatedTime = time.Now()
	db.Sex = info.Sex
	db.Name = info.RealName
	db.Job = info.Job
	db.Phone = info.Phone
	err := nosql.CreateDatum(db)
	if err == nil {
		info.initInfo(db)
	}
	return err
}

func AllUsers() []*UserInfo {
	return cacheCtx.users
}

func GetUser(uid string) *UserInfo {
	for i := 0;i < len(cacheCtx.users);i += 1 {
		if cacheCtx.users[i].UID == uid {
			return cacheCtx.users[i]
		}
	}
	db,err := nosql.GetUser(uid)
	if err == nil {
		user := new(UserInfo)
		user.initInfo(db)
		cacheCtx.users = append(cacheCtx.users, user)
		return user
	}
	return nil
}

func GetUserByAccount(account string) *UserInfo {
	for i := 0;i < len(cacheCtx.users);i += 1 {
		if cacheCtx.users[i].Account == account {
			return cacheCtx.users[i]
		}
	}
	db,err := nosql.GetUserByAccount(account)
	if err == nil {
		user := new(UserInfo)
		user.initInfo(db)
		cacheCtx.users = append(cacheCtx.users, user)
		return user
	}
	return nil
}

func RemoveUser(uid, operator string) error {
	if len(uid) < 1{
		return errors.New("the user uid is empty")
	}
	err := nosql.RemoveUser(uid, operator)
	if err == nil {
		for i := 0;i < len(cacheCtx.users);i += 1 {
			if cacheCtx.users[i].UID == uid {
				cacheCtx.users = append(cacheCtx.users[:i], cacheCtx.users[i+1:]...)
				break
			}
		}
	}
	return err
}

func (mine *UserInfo)initInfo(db *nosql.User)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Type = db.Type
	mine.Operator = db.Operator
	mine.Creator = db.Creator
	datum,err := nosql.GetDatum(db.Datum)
	if err == nil {
		mine.Datum.initInfo(datum)
	}
}

func (mine *UserInfo)UpdateBase(name, real, phone, remark, job, operator string, sex uint8) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) <1 {
		remark = mine.Remark
	}
	err := nosql.UpdateUserBase(mine.UID, name, remark, operator)
	if err != nil {
		return err
	}
	if len(real) < 1 {
		real = mine.Datum.RealName
	}

	if len(phone) < 1 {
		phone = mine.Datum.Phone
	}

	if len(job) < 1 {
		job = mine.Datum.Job
	}

	err1 := nosql.UpdateDatumBase(mine.Datum.UID, real, phone, job, sex)
	if err1 != nil{
		return err1
	}
	mine.Name = name
	mine.Datum.RealName = real
	mine.Remark = remark
	mine.Datum.Phone = phone
	mine.Datum.Sex = sex
	mine.Datum.Job = job
	mine.Operator = operator
	mine.UpdateTime = time.Now()
	return  nil
}

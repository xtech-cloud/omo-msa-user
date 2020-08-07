package cache

import (
	"omo.msa.user/proxy/nosql"
	"time"
)

type UserInfo struct {
	Type uint8
	BaseInfo
	NickName string
	Account string
	Remark string
	Phone string
	Sex uint8
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
	mine.Sex = db.Sex
	mine.Phone = db.Phone
	mine.NickName = db.Nick
}

func (mine *UserInfo)UpdateBase(name, nick, remark, operator string, sex uint8) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) <1 {
		remark = mine.Remark
	}
	if len(nick) < 1 {
		nick = mine.NickName
	}
	err := nosql.UpdateUserBase(mine.UID, name, nick, remark, operator, sex)
	if err == nil {
		mine.Name = name
		mine.Sex = sex
		mine.Remark = remark
		mine.NickName = nick
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}

	return  nil
}

func (mine *UserInfo)UpdatePhone(phone, operator string) error {
	err := nosql.UpdateUserPhone(mine.UID, phone, operator)
	if err == nil {
		mine.Phone = phone
		mine.Operator = operator
	}
	return err
}

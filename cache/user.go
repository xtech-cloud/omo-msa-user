package cache

import (
	"errors"
	"omo.msa.user/proxy"
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
	Entity string
	SNS []proxy.SNSInfo
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
	mine.Account = db.Account
	mine.Entity = db.Entity
	mine.SNS = db.SNS
	if mine.SNS == nil {
		mine.SNS = make([]proxy.SNSInfo, 0, 1)
	}
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

func (mine *UserInfo)UpdateEntity(entity, operator string) error {
	if operator == ""{
		operator = mine.Operator
	}
	err := nosql.UpdateUserEntity(mine.UID, entity, operator)
	if err == nil {
		mine.Entity = entity
		mine.Operator = operator
	}
	return err
}

func (mine *UserInfo)HadSNS(uid string) bool {
	for _, sn := range mine.SNS {
		if sn.UID == uid {
			return true
		}
	}
	return false
}

func (mine *UserInfo)AppendSNS(uid, name string, kind uint8) error {
	if uid == "" {
		return errors.New("the sns uid is empty")
	}
	if mine.HadSNS(uid) {
		return nil
	}
	tmp := proxy.SNSInfo{UID: uid, Name: name, Type: kind}
	err := nosql.AppendUserSNS(mine.UID, tmp)
	if err == nil {
		mine.SNS = append(mine.SNS, tmp)
	}
	return err
}

func (mine *UserInfo)SubtractSNS(uid string) error {
	if uid == "" {
		return errors.New("the sns uid is empty")
	}
	if !mine.HadSNS(uid) {
		return nil
	}
	err := nosql.SubtractUserSNS(mine.UID, uid)
	if err == nil {
		for i := 0; i < len(mine.SNS);i += 1 {
			if mine.SNS[i].UID == uid {
				mine.SNS = append(mine.SNS[:i], mine.SNS[i:]...)
				break
			}
		}
	}
	return err
}

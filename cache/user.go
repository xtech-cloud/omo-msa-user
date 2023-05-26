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
	Account  string
	Remark   string
	Phone    string
	Sex      uint8
	Status   uint8
	Entity   string
	Portrait string
	Tags     []string
	Follows  []string
	Relates  []string
	SNS      []proxy.SNSInfo
}

func (mine *UserInfo) initInfo(db *nosql.User, st uint8) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.DeleteTime = db.DeleteTime
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Type = db.Type
	mine.Operator = db.Operator
	mine.Creator = db.Creator
	mine.Sex = db.Sex
	mine.Phone = db.Phone
	mine.Status = st
	mine.NickName = db.Nick
	mine.Account = db.Account
	mine.Entity = db.Entity
	mine.Portrait = db.Portrait
	mine.Tags = db.Tags
	if mine.Tags == nil {
		mine.Tags = make([]string, 0, 5)
	}
	mine.Follows = db.Follows
	if mine.Follows == nil {
		mine.Follows = make([]string, 0, 5)
	}
	mine.Relates = db.Relates
	if mine.Relates == nil {
		mine.Relates = make([]string, 0, 5)
	}
	mine.SNS = db.SNS
	if mine.SNS == nil {
		mine.SNS = make([]proxy.SNSInfo, 0, 1)
	}
}

func (mine *UserInfo) UpdateBase(name, nick, remark, portrait, operator string, sex uint8) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	if len(nick) < 1 {
		nick = mine.NickName
	}
	if len(portrait) < 1 {
		portrait = mine.Portrait
	}
	err := nosql.UpdateUserBase(mine.UID, name, nick, remark, portrait, operator, sex)
	if err == nil {
		mine.Name = name
		mine.Sex = sex
		mine.Remark = remark
		mine.NickName = nick
		mine.Operator = operator
		mine.Portrait = portrait
		mine.UpdateTime = time.Now()
	}

	return nil
}

func (mine *UserInfo) UpdateFollows(list []string) error {
	err := nosql.UpdateUserFollows(mine.UID, list)
	if err == nil {
		mine.Follows = list
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *UserInfo) UpdateRelates(list []string) error {
	err := nosql.UpdateUserRelates(mine.UID, list)
	if err == nil {
		mine.Relates = list
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *UserInfo) UpdatePortrait(icon, operator string) error {
	if len(icon) < 2 {
		return nil
	}
	err := nosql.UpdateUserPortrait(mine.UID, icon, operator)
	if err == nil {
		mine.Portrait = icon
		mine.Operator = operator
	}
	return err
}

func (mine *UserInfo) UpdatePhone(phone, operator string) error {
	if len(phone) < 7 {
		return errors.New("the phone format is error")
	}
	err := nosql.UpdateUserPhone(mine.UID, phone, operator)
	if err == nil {
		mine.Phone = phone
		mine.Operator = operator
	}
	return err
}

func (mine *UserInfo) UpdateType(kind uint8) error {
	if kind < 1 {
		return errors.New("the user type is error")
	}
	err := nosql.UpdateUserType(mine.UID, kind)
	if err == nil {
		mine.Type = kind
	}
	return err
}

func (mine *UserInfo) UpdateEntity(entity, operator string) error {
	if entity == "" {
		return nil
	}
	if operator == "" {
		operator = mine.Operator
	}
	err := nosql.UpdateUserEntity(mine.UID, entity, operator)
	if err == nil {
		mine.Entity = entity
		mine.Operator = operator
	}
	return err
}

func (mine *UserInfo) UpdateTags(operator string, tags []string) error {
	if tags == nil {
		return nil
	}
	if operator == "" {
		operator = mine.Operator
	}
	err := nosql.UpdateUserTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *UserInfo) HadSNS(uid string) bool {
	for _, sn := range mine.SNS {
		if sn.UID == uid {
			return true
		}
	}
	return false
}

func (mine *UserInfo) HadSNSByType(kind uint8) bool {
	for _, sn := range mine.SNS {
		if sn.Type == kind {
			return true
		}
	}
	return false
}

func (mine *UserInfo) AppendSNS(uid, name string, kind uint8) error {
	if len(uid) < 3 {
		return errors.New("the sns uid is empty")
	}
	if mine.HadSNS(uid) {
		return nil
	}
	//if mine.HadSNSByType(kind) {
	//	return errors.New("the sns type is exist")
	//}
	tmp := proxy.SNSInfo{UID: uid, Name: name, Type: kind}
	err := nosql.AppendUserSNS(mine.UID, tmp)
	if err == nil {
		mine.SNS = append(mine.SNS, tmp)
	}
	return err
}

func (mine *UserInfo) SubtractSNS(uid string) error {
	if uid == "" {
		return errors.New("the sns uid is empty")
	}
	if !mine.HadSNS(uid) {
		return nil
	}
	err := nosql.SubtractUserSNS(mine.UID, uid)
	if err == nil {
		for i := 0; i < len(mine.SNS); i += 1 {
			if mine.SNS[i].UID == uid {
				mine.SNS = append(mine.SNS[:i], mine.SNS[i+1:]...)
				break
			}
		}
	}
	return err
}

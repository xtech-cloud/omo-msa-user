package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.user/proxy/nosql"
	"time"
)

type WechatInfo struct {
	BaseInfo
	Sex       uint8
	OpenID    string
	UnionID   string
	Code      string
	Portrait  string
}

func (mine *cacheContext)AddWechat(info *WechatInfo) (*WechatInfo, error) {
	if info == nil {
		return nil,errors.New("the wechat info is nil")
	}
	return mine.CreateWechat(info.Name, info.OpenID, info.UnionID, info.Portrait, info.Creator)
}

func (mine *cacheContext)CreateWechat(name, open, union, img, creator string) (*WechatInfo, error) {
	tmp := mine.GetWechatByOpen(open)
	if tmp != nil {
		return tmp,nil
	}
	var db = new(nosql.Wechat)
	db.UID = primitive.NewObjectID()
	db.CreatedTime = time.Now()
	db.ID = nosql.GetWechatNextID()
	db.Creator = creator
	db.Name = name
	db.OpenID = open
	db.UUID = union
	db.Sex = 0
	db.Portrait = img
	err := nosql.CreateWechat(db)
	if err != nil {
		return nil, err
	}
	info := new(WechatInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext)RemoveWechat(uid string) bool {
	//length := len(mine.wechats)
	//for i := 0; i < length; i++ {
	//	if mine.wechats[i].UID == uid {
	//		mine.wechats = append(mine.wechats[:i], mine.wechats[i+1:]...)
	//		return true
	//	}
	//}
	//return false
	return true
}

func (mine *cacheContext)GetWechatByOpen(uid string) *WechatInfo {
	//for i := 0; i < len(mine.wechats); i += 1 {
	//	if mine.wechats[i].OpenID == uid {
	//		return mine.wechats[i]
	//	}
	//}
	db,_ := nosql.GetWechatByOpen(uid)
	if db != nil {
		var info = new(WechatInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext)GetWechat(uid string) *WechatInfo {
	//for i := 0; i < len(mine.wechats); i += 1 {
	//	if mine.wechats[i].UID == uid {
	//		return mine.wechats[i]
	//	}
	//}
	db,_ := nosql.GetWechat(uid)
	if db != nil {
		var info = new(WechatInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *WechatInfo) initInfo(db *nosql.Wechat) bool {
	if db == nil {
		return false
	}
	mine.UID = db.UID.Hex()
	mine.OpenID = db.OpenID
	mine.UnionID = db.UUID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.ID = db.ID
	mine.Sex = db.Sex
	mine.Portrait = db.Portrait
	return true
}

func (mine *WechatInfo)UpdateBase(name, open, union, img, operator string) error {
	if name == ""{
		name = mine.Name
	}
	if open == ""{
		open = mine.OpenID
	}
	if union == "" {
		union = mine.UnionID
	}
	if img == "" {
		img = mine.Portrait
	}
	err := nosql.UpdateWechatBase(mine.UID, name, open, union, img, operator)
	if err == nil {
		mine.Name = name
		mine.OpenID = open
		mine.UnionID = union
		mine.Portrait = img
		mine.Operator = operator
	}
	return err
}

package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.user/proxy/nosql"
	"omo.msa.user/tool"
	"time"
)

const (
	MessageSleep MessageStatus = 0
	MessageAwake MessageStatus = 1
	MessageRead  MessageStatus = 2
)

type MessageStatus uint8

type MessageInfo struct {
	Type   uint8 //消息类型，活动，通知
	Status MessageStatus
	Stamp  uint64 //生效时间戳
	BaseInfo
	Owner   string   //创建者所属场景
	User    string   //用户
	Quote   string   //活动，通知uid等
	Targets []string //用户下面的目标实体
}

func (mine *MessageInfo) Awake() error {
	if mine.Status == MessageAwake {
		return nil
	}
	err := nosql.UpdateMessageStatus(mine.UID, uint8(MessageAwake))
	if err == nil {
		mine.Status = MessageAwake
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *MessageInfo) Read() error {
	if mine.Status == MessageRead {
		return nil
	}
	err := nosql.UpdateMessageStatus(mine.UID, uint8(MessageRead))
	if err == nil {
		mine.Status = MessageRead
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *MessageInfo) UpdateTargets(arr []string) error {
	list := make([]string, 0, len(arr)+len(mine.Targets))
	list = append(list, mine.Targets...)
	for _, s := range arr {
		if !tool.HasItem(list, s) {
			list = append(list, s)
		}
	}
	err := nosql.UpdateMessageTargets(mine.UID, list)
	if err == nil {
		mine.Targets = list
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *MessageInfo) initInfo(db *nosql.Message) {
	mine.UID = db.UID.Hex()
	mine.Creator = db.Creator
	mine.CreateTime = db.CreatedTime
	mine.Name = db.Name
	mine.Stamp = db.Stamp
	mine.Status = MessageStatus(db.Status)
	mine.Owner = db.Owner
	mine.Quote = db.Quote
	mine.Targets = db.Targets
	mine.User = db.User
	mine.Type = db.Type
}

func (mine *cacheContext) CreateMessage(owner, user, quote, operator string, tp uint32, stamp uint64, targets []string) error {
	db := new(nosql.Message)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetMessageNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.UpdatedTime = time.Now()
	db.User = user
	db.Targets = targets
	db.Type = uint8(tp)
	db.Status = uint8(MessageSleep)
	db.Owner = owner
	db.Quote = quote
	db.Stamp = stamp

	return nosql.CreateMessage(db)
}

func (mine *cacheContext) GetMessagesByUser(user string) []*MessageInfo {
	if user == "" {
		return nil
	}
	dbs, _ := nosql.GetMessagesByUser(user)
	list := make([]*MessageInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(MessageInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetMessagesByQuote(user, quote string) (*MessageInfo, error) {
	if user == "" || quote == "" {
		return nil, errors.New("the user or quote is empty")
	}
	db, er := nosql.GetMessagesByQuote(user, quote)
	if er == nil {
		info := new(MessageInfo)
		info.initInfo(db)
		return info, nil
	}

	return nil, er
}

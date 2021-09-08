package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.user/proxy"
	"time"
)

type User struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name    string `json:"name" bson:"name"`
	Account string `json:"account" bson:"account"`
	Remark  string `json:"remark" bson:"remark"`
	Type    uint8  `json:"type" bson:"type"`
	Nick    string `json:"nick" bson:"nick"`
	Phone   string `json:"phone" bson:"phone"`
	Sex     uint8  `json:"sex" bson:"sex"`
	Entity  string `json:"entity" bson:"entity"`
	Portrait string `json:"portrait" bson:"portrait"`
	Tags  []string `json:"tags" bson:"tags"`
	SNS     []proxy.SNSInfo `json:"sns" bson:"sns"`
	Follows []string `json:"follows" bson:"follows"`
}

func CreateUser(info *User) error {
	_, err := insertOne(TableUser, info)
	if err != nil {
		return err
	}
	return nil
}

func GetUserNextID() uint64 {
	num, _ := getSequenceNext(TableUser)
	return num
}

func GetUser(uid string) (*User, error) {
	result, err := findOne(TableUser, uid)
	if err != nil {
		return nil, err
	}
	model := new(User)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetUserByID(id uint64) (*User, error) {
	msg := bson.M{"id": id}
	result, err := findOneBy(TableUser, msg)
	if err != nil {
		return nil, err
	}
	model := new(User)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllUsers() ([]*User, error) {
	cursor, err1 := findAll(TableUser, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*User, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(User)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetUserCount() int64 {
	num, _ := getCount(TableUser)
	return num
}

func GetUsersByAccount(uid string) ([]*User, error) {
	msg := bson.M{"account": uid, "deleteAt":new(time.Time)}
	cursor, err := findMany(TableUser, msg, 0)
	if err != nil {
		return nil, err
	}
	var items = make([]*User, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(User)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetUsersByType(kind uint8) ([]*User, error) {
	msg := bson.M{"type": kind, "deleteAt":new(time.Time)}
	cursor, err := findMany(TableUser, msg, 0)
	if err != nil {
		return nil, err
	}
	var items = make([]*User, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(User)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetUserByPhone(phone string) (*User, error) {
	msg := bson.M{"phone": phone, "deleteAt":new(time.Time)}
	result, err := findOneBy(TableUser, msg)
	if err != nil {
		return nil, err
	}
	model := new(User)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetUserByEntity(entity string) (*User, error) {
	msg := bson.M{"entity": entity, "deleteAt":new(time.Time)}
	result, err := findOneBy(TableUser, msg)
	if err != nil {
		return nil, err
	}
	model := new(User)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetUserBySNS(uid string) (*User, error) {
	msg := bson.M{"sns.uid": uid, "deleteAt":new(time.Time)}
	result, err := findOneBy(TableUser, msg)
	if err != nil {
		return nil, err
	}
	model := new(User)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateUserBase(uid, name, nick, remark, portrait, operator string, sex uint8) error {
	msg := bson.M{"name": name, "nick": nick, "remark": remark, "portrait": portrait, "sex": sex, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserPhone(uid, phone, operator string) error {
	msg := bson.M{"phone": phone, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserType(uid string, kind uint8) error {
	msg := bson.M{"type": kind, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserEntity(uid, entity, operator string) error {
	msg := bson.M{"entity": entity, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserPortrait(uid string, icon, operator string) error {
	msg := bson.M{"portrait": icon, "operator": operator,  "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserFollows(uid string, list []string) error {
	msg := bson.M{"follows": list, "operator": uid,  "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func RemoveUser(uid, operator string) error {
	_, err := removeOne(TableUser, uid, operator)
	return err
}

func AppendUserSNS(uid string, sns proxy.SNSInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"sns": sns}
	_, err := appendElement(TableUser, uid, msg)
	return err
}

func SubtractUserSNS(uid, sns string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"sns": bson.M{"uid": sns}}
	_, err := removeElement(TableUser, uid, msg)
	return err
}
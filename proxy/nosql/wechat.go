package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Wechat struct {
	Sex         uint8              `json:"sex" bson:"sex"`
	UID         primitive.ObjectID `bson:"_id"`
	OpenID      string             `json:"oid" bson:"oid"`
	UUID        string             `json:"uuid" bson:"uuid"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Creator     string `json:"creator" bson:"creator"`
	Operator    string `json:"operator" bson:"operator"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Portrait    string             `json:"port" bson:"port"`
}

func CreateWechat(info *Wechat) error {
	_, err := insertOne(TableWechat, info)
	if err != nil {
		return err
	}
	return nil
}

func GetWechatNextID() uint64 {
	num, _ := getSequenceNext(TableWechat)
	return num
}

func GetWechat(uid string) (*Wechat, error) {
	result, err := findOne(TableWechat, uid)
	if err != nil {
		return nil, err
	}
	model := new(Wechat)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateWechatBase(uid string, name, open, union, port, operator string) error {
	msg := bson.M{"name": name, "oid": open, "uuid": union, "port": port,"operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableWechat, uid, msg)
	return err
}

func GetWechatByOpen(uid string) (*Wechat, error) {
	msg := bson.M{"oid": uid}
	result, err := findOneBy(TableWechat, msg)
	if err != nil {
		return nil, err
	}
	model := new(Wechat)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllWechats() ([]*Wechat, error) {
	cursor, err1 := findAll(TableWechat, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Wechat, 0, 2000)
	for cursor.Next(context.Background()) {
		var node = new(Wechat)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveWechat(uid, operator string) error {
	_, err := removeOne(TableWechat, uid, operator)
	return err
}

func dropWechat() error {
	err := dropOne(TableWechat)
	return err
}

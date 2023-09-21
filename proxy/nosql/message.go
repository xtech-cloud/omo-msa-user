package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name string `json:"name" bson:"name"`
	// 目标类型
	Type    uint8    `json:"type" bson:"type"`
	Stamp   uint64   `json:"stamp" bson:"stamp"`
	User    string   `json:"user" bson:"user"`
	Owner   string   `json:"owner" bson:"owner"`
	Status  uint8    `json:"status" bson:"status"`
	Quote   string   `json:"quote" bson:"quote"`
	Targets []string `json:"targets" bson:"targets"`
}

func CreateMessage(info *Message) error {
	_, err := insertOne(TableMessage, info)
	if err != nil {
		return err
	}
	return nil
}

func GetMessageNextID() uint64 {
	num, _ := getSequenceNext(TableMessage)
	return num
}

func GetMessage(uid string) (*Message, error) {
	result, err := findOne(TableMessage, uid)
	if err != nil {
		return nil, err
	}
	model := new(Message)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetMessageHistories(user string, kind uint8) ([]*Message, error) {
	msg := bson.M{"user": user, "type": kind, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableMessage, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Message, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Message)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMessageByType(kind uint8) ([]*Message, error) {
	msg := bson.M{"type": kind, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableMessage, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Message, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Message)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMessageByTarget(target string) ([]*Message, error) {
	msg := bson.M{"targets": target, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableMessage, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Message, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Message)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMessagesByUser(user string) ([]*Message, error) {
	msg := bson.M{"user": user, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableMessage, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Message, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Message)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMessagesByQuote(user, quote string) (*Message, error) {
	msg := bson.M{"user": user, "quote": quote, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableMessage, msg)
	if err != nil {
		return nil, err
	}
	model := new(Message)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetMessageCountByUser(user string) (int64, error) {
	msg := bson.M{"user": user, "deleteAt": new(time.Time)}
	return getCountByFilter(TableMessage, msg)
}

func UpdateMessageStatus(uid string, st uint8) error {
	msg := bson.M{"status": st, "updatedAt": time.Now()}
	_, err := updateOne(TableMessage, uid, msg)
	return err
}

func UpdateMessageTargets(uid string, targets []string) error {
	msg := bson.M{"targets": targets, "updatedAt": time.Now()}
	_, err := updateOne(TableMessage, uid, msg)
	return err
}

func GetAllMessages() ([]*Message, error) {
	cursor, err1 := findAll(TableMessage, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Message, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Message)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveMessage(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	_, err := removeOne(TableMessage, uid, operator)
	return err
}

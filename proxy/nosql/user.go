package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Datum   string `json:"datum" bson:"datum"`
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

func GetUserByAccount(uid string) (*User, error) {
	msg := bson.M{"account": uid}
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

func UpdateUserBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserCover(uid string, icon string) error {
	msg := bson.M{"cover": icon, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func RemoveUser(uid, operator string) error {
	_, err := removeOne(TableUser, uid, operator)
	return err
}

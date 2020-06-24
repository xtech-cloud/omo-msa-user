package nosql

import (
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

	Name   string                `json:"name" bson:"name"`
	Account  string                `json:"table" bson:"table"`
	Remark string `json:"remark" bson:"remark"`
	Datum string                `json:"datum" bson:"datum"`
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

func GetUserByAccount(uid string) (*User, error) {
	msg := bson.M{"account":uid}
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

func UpdateUserBase(uid, name, remark string) error {
	msg := bson.M{"name": name, "remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserCover(uid string, icon string) error {
	msg := bson.M{"cover": icon, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func RemoveUser(uid string) error {
	_, err := removeOne(TableUser, uid)
	return err
}

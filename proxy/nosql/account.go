package nosql


import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Account struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name    string `json:"name" bson:"name"`
	Status  uint8 `json:"status" bson:"status"`
	Passwords string `json:"passwords" bson:"passwords"`
}

func CreateAccount(info *Account) error {
	_, err := insertOne(TableAccount, info)
	if err != nil {
		return err
	}
	return nil
}

func GetAccountNextID() uint64 {
	num, _ := getSequenceNext(TableAccount)
	return num
}

func GetAccountCount() int64 {
	num, _ := getCount(TableAccount)
	return num
}

func GetAccount(uid string) (*Account, error) {
	result, err := findOne(TableAccount, uid)
	if err != nil {
		return nil, err
	}
	model := new(Account)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAccountByName(name string) (*Account, error) {
	msg := bson.M{"name": name}
	result, err := findOneBy(TableAccount, msg)
	if err != nil {
		return nil, err
	}
	model := new(Account)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func getOldAccounts() ([]*Account, error) {
	var items = make([]*Account, 0, 20)
	filter := bson.M{"status": bson.M{"$exists": false}}
	cursor, err1 := findMany(TableAccount, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Account)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}


func GetAllAccounts() ([]*Account, error) {
	cursor, err1 := findAll(TableAccount, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Account, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(Account)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateAccountBase(uid, name, operator string) error {
	msg := bson.M{"name": name, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAccount, uid, msg)
	return err
}

func UpdateAccountStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAccount, uid, msg)
	return err
}

func UpdateAccountPasswords(uid, psw, operator string) error {
	msg := bson.M{"passwords": psw, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAccount, uid, msg)
	return err
}

func RemoveAccount(uid, operator string) error {
	_, err := removeOne(TableAccount, uid, operator)
	return err
}

func DeleteAccount(uid string) error {
	_, err := deleteOne(TableAccount, uid)
	return err
}
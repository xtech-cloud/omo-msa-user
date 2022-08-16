package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Behaviour struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	// 目标类型
	Type uint8    `json:"type" bson:"type"`
	User   string   `json:"user" bson:"user"`
	Target string   `json:"target" bson:"target"`
	// 动作
	Action   uint8 `json:"action" bson:"action"`
}

func CreateBehaviour(info *Behaviour) error {
	_, err := insertOne(TableBehaviour, info)
	if err != nil {
		return err
	}
	return nil
}

func GetBehaviourNextID() uint64 {
	num, _ := getSequenceNext(TableBehaviour)
	return num
}

func GetBehaviour(uid string) (*Behaviour, error) {
	result, err := findOne(TableBehaviour, uid)
	if err != nil {
		return nil, err
	}
	model := new(Behaviour)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetBehaviourHistories(user string, act, kind uint8, num int64) ([]*Behaviour, error) {
	msg := bson.M{"user": user, "action": act, "type":kind}
	cursor, err1 := findMany(TableBehaviour, msg, num)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Behaviour, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Behaviour)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetBehaviourByType(user string, kind uint8, num int64) ([]*Behaviour, error) {
	msg := bson.M{"user": user, "type":kind}
	cursor, err1 := findMany(TableBehaviour, msg, num)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Behaviour, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Behaviour)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetBehaviourByAct(user, target string,kind uint8) (*Behaviour, error) {
	msg := bson.M{"user": user, "target": target, "type":kind, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableBehaviour, msg)
	if err != nil {
		return nil, err
	}
	model := new(Behaviour)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetBehaviourByTarget(user, target string) (*Behaviour, error) {
	msg := bson.M{"user": user, "target": target, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableBehaviour, msg)
	if err != nil {
		return nil, err
	}
	model := new(Behaviour)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetBehavioursByTarget(user, target string) ([]*Behaviour, error) {
	msg := bson.M{"user": user, "target": target, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableBehaviour, msg, 50)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Behaviour, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Behaviour)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetBehavioursByAction(target string, act uint8) ([]*Behaviour, error) {
	msg := bson.M{"target": target, "action":act, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableBehaviour, msg, 50)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Behaviour, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Behaviour)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetBehaviourCountByAction(target string, act uint8) (int64, error) {
	msg := bson.M{"target": target, "action":act, "deleteAt": new(time.Time)}
	return getCountByFilter(TableBehaviour, msg)
}

func GetBehaviourCountByUser(user string) (int64, error) {
	msg := bson.M{"user": user, "deleteAt": new(time.Time)}
	return getCountByFilter(TableBehaviour, msg)
}

func UpdateBehaviourAction(user, target string, act uint8) error {
	filter := bson.M{"user":user, "target": target}
	update := bson.M{"$set": bson.M{"action": act, "updatedAt": time.Now()}}
	_, err := updateOneBy(TableBehaviour, filter, update)
	return err
}

func GetAllBehaviours() ([]*Behaviour, error) {
	cursor, err1 := findAll(TableBehaviour, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Behaviour, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Behaviour)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveBehaviour(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	_, err := removeOne(TableBehaviour, uid, operator)
	return err
}

package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"omo.msa.user/proxy"
	"time"
)

type Score struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Type   uint8             `json:"type" bson:"type"`
	Count  uint64            `json:"count" bson:"count"`
	Stamp  int64             `json:"stamp" bson:"stamp"` //以天为单位
	Scene  string            `json:"scene" bson:"scene"`
	Entity string            `json:"entity" bson:"entity"`
	List   []proxy.ScorePair `json:"list" bson:"list"`
}

func CreateScore(info *Score) error {
	_, err := insertOne(TableScores, &info)
	return err
}

func GetScoreNextID() uint64 {
	num, _ := getSequenceNext(TableScores)
	return num
}

func GetScoreCount() int64 {
	num, _ := getCount(TableScores)
	return num
}

func RemoveScore(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db score uid is empty ")
	}
	_, err := removeOne(TableScores, uid, operator)
	return err
}

func GetScore(uid string) (*Score, error) {
	if len(uid) < 2 {
		return nil, errors.New("db score uid is empty of GetScore")
	}

	result, err := findOne(TableScores, uid)
	if err != nil {
		return nil, err
	}
	model := new(Score)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetScoreBySceneEntity(scene, entity string) (*Score, error) {
	if len(scene) < 2 || len(entity) < 2 {
		return nil, errors.New("db score scene or entity is empty of GetScoreBySceneEntity")
	}
	filter := bson.M{"scene": scene, "entity": entity}
	result, err := findOneBy(TableScores, filter)
	if err != nil {
		return nil, err
	}
	model := new(Score)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetScoreBySceneDate(scene string, date int64) (*Score, error) {
	if len(scene) < 2 || date < 2 {
		return nil, errors.New("db score scene or stamp is empty of GetScoreBySceneDate")
	}
	filter := bson.M{"scene": scene, "stamp": date}
	result, err := findOneBy(TableScores, filter)
	if err != nil {
		return nil, err
	}
	model := new(Score)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetScoreBySceneDur(scene string, from, to int64) ([]*Score, error) {
	if len(scene) < 2 {
		return nil, errors.New("db score scene or stamp from or to is empty of GetScoreBySceneDur")
	}
	filter := bson.M{"scene": scene, "stamp": bson.M{"$gte": from, "$lte": to}}
	cursor, err1 := findMany(TableScores, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Score, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Score)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetScoresByTop(num int64, tp uint8) ([]*Score, error) {
	var items = make([]*Score, 0, 20)
	filter := bson.M{"type": tp}
	opts := options.Find().SetSort(bson.D{{"count", -1}}).SetLimit(num)
	cursor, err1 := findManyByOpts(TableScores, filter, opts)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Score)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetScoresBySceneTop(scene string, num int64, tp uint8) ([]*Score, error) {
	var items = make([]*Score, 0, 20)
	filter := bson.M{"scene": scene, "type": tp}
	opts := options.Find().SetSort(bson.D{{"count", -1}}).SetLimit(num)
	cursor, err1 := findManyByOpts(TableScores, filter, opts)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Score)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateScorePairs(uid, operator string, count uint64, arr []proxy.ScorePair) error {
	if len(uid) < 2 {
		return errors.New("db score uid is empty of GetAsset")
	}

	msg := bson.M{"list": arr, "count": count, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableScores, uid, msg)
	return err
}

func UpdateScoreCount(uid, operator string, count uint64) error {
	if len(uid) < 2 {
		return errors.New("db score uid is empty of GetAsset")
	}

	msg := bson.M{"count": count, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableScores, uid, msg)
	return err
}

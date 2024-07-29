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

type User struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name     string          `json:"name" bson:"name"`
	Account  string          `json:"account" bson:"account"`
	Remark   string          `json:"remark" bson:"remark"`
	Type     uint8           `json:"type" bson:"type"`
	Nick     string          `json:"nick" bson:"nick"`
	Phone    string          `json:"phone" bson:"phone"`
	Sex      uint8           `json:"sex" bson:"sex"`
	Entity   string          `json:"entity" bson:"entity"`
	Portrait string          `json:"portrait" bson:"portrait"`
	Shown    proxy.ShownInfo `json:"shown" bson:"shown"`
	Tags     []string        `json:"tags" bson:"tags"`
	Follows  []string        `json:"follows" bson:"follows"`
	Relates  []string        `json:"relates" bson:"relates"`
	SNS      []proxy.SNSInfo `json:"sns" bson:"sns"`
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
	msg := bson.M{"account": uid}
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
	msg := bson.M{"type": kind}
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
	msg := bson.M{"phone": phone}
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
	msg := bson.M{"entity": entity}
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
	msg := bson.M{"sns.uid": uid}
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

func GetUsersByScenePage(scene string, start, num int64) ([]*User, error) {
	def := new(time.Time)
	filter := bson.M{"relates": scene, "deleteAt": def}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}}).SetLimit(num).SetSkip(start)
	cursor, err1 := findManyByOpts(TableUser, filter, opts)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*User, 0, 20)
	for cursor.Next(context.TODO()) {
		var node = new(User)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetUsersByPage(start, num int64) ([]*User, error) {
	def := new(time.Time)
	filter := bson.M{"deleteAt": def}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}}).SetLimit(num).SetSkip(start)
	cursor, err1 := findManyByOpts(TableUser, filter, opts)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*User, 0, 20)
	for cursor.Next(context.TODO()) {
		var node = new(User)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetUsersCount() int64 {
	num, err1 := getCount(TableUser)
	if err1 != nil {
		return num
	}

	return num
}

func GetUsersCountByScene(scene string) int64 {
	def := new(time.Time)
	filter := bson.M{"relates": scene, "deleteAt": def}
	num, err1 := getCountByFilter(TableUser, filter)
	if err1 != nil {
		return num
	}

	return num
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
	msg := bson.M{"portrait": icon, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserFollows(uid string, list []string) error {
	msg := bson.M{"follows": list, "operator": uid, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserRelates(uid string, list []string) error {
	msg := bson.M{"relates": list, "operator": uid, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserShown(uid string, shown proxy.ShownInfo) error {
	msg := bson.M{"shown": shown, "operator": uid, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func DeleteUser(uid string) error {
	_, err := deleteOne(TableUser, uid)
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

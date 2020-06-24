package nosql

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Datum struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Name   string `json:"name" bson:"name"`  //真名
	Phone  string `json:"phone" bson:"phone"`
	Job    string `json:"job" bson:"job"`
	Sex    uint8  `json:"sex" bson:"sex"`
}

func CreateDatum(info *Datum) error {
	_, err := insertOne(TableDatum, info)
	if err != nil {
		return err
	}
	return nil
}

func GetDatumNextID() uint64 {
	num, _ := getSequenceNext(TableDatum)
	return num
}

func GetDatum(uid string) (*Datum, error) {
	result, err := findOne(TableDatum, uid)
	if err != nil {
		return nil, err
	}
	model := new(Datum)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateDatumBase(uid, name, phone, job string, sex uint8) error {
	msg := bson.M{"name": name, "phone": phone, "job":job, "sex": sex, "updatedAt": time.Now()}
	_, err := updateOne(TableDatum, uid, msg)
	return err
}

func UpdateDatumCover(uid string, icon string) error {
	msg := bson.M{"cover": icon, "updatedAt": time.Now()}
	_, err := updateOne(TableDatum, uid, msg)
	return err
}

func RemoveDatum(uid string) error {
	_, err := removeOne(TableDatum, uid)
	return err
}

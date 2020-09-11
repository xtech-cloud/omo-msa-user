package proxy

type SNSInfo struct {
	Type uint8 `json:"type" bson:"type"`
	Name string `json:"name" bson:"name"`
	UID string `json:"uid" bson:"uid"`
}
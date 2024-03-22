package proxy

type SNSInfo struct {
	Type uint8  `json:"type" bson:"type"`
	Name string `json:"name" bson:"name"`
	UID  string `json:"uid" bson:"uid"`
	ID   string `json:"id" bson:"id"`
}

type ScorePair struct {
	Type  uint32 `json:"type" bson:"type"`
	Count uint32 `json:"count" bson:"count"`
}

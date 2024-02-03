package cache

import (
	"github.com/micro/go-micro/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.user/proxy/nosql"
	"strings"
	"time"
)

const (
	ActionUnknown ActionType = 0
	ActionVisit   ActionType = 1
	ActionCare    ActionType = 2
)

const (
	TargetTypeAlbum      TargetType = 1
	TargetTypeCollective TargetType = 2
	TargetTypeActivity   TargetType = 3
)

type ActionType uint8

type TargetType uint8

func (mine *cacheContext) createBehaviour(user, target string, kind TargetType, act ActionType) error {
	db := new(nosql.Behaviour)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetBehaviourNextID()
	db.CreatedTime = time.Now()
	db.Creator = user
	db.UpdatedTime = time.Now()
	db.User = user
	db.Target = target
	db.Type = uint8(kind)
	db.Action = uint8(act)

	err := nosql.CreateBehaviour(db)
	if err != nil {
		logger.Error("create behaviour failed that err = " + err.Error())
	}
	return err
}

func (mine *cacheContext) removeBehaviour(user, target string) error {
	db, err := nosql.GetBehaviourByTarget(user, target)
	if err != nil {
		return err
	}
	return nosql.RemoveBehaviour(db.UID.Hex(), "")
}

func (mine *cacheContext) UpdateBehaviour(user, target string, act ActionType) error {
	err := nosql.UpdateBehaviourAction(user, target, uint8(act))
	return err
}

func (mine *cacheContext) AddBehaviour(user, target string, kind TargetType, act ActionType) error {
	had, err := mine.HadBehaviour(user, target)
	if err != nil {
		return err
	}
	if had {
		return nil
	}
	return mine.createBehaviour(user, target, kind, act)
}

func (mine *cacheContext) HadBehaviour(user, target string) (bool, error) {
	db, err := nosql.GetBehaviourByTarget(user, target)
	if err != nil && !strings.Contains(err.Error(), "no documents in result") {
		return false, err
	}
	if db != nil {
		return true, nil
	}
	return false, nil
}

func (mine *cacheContext) HadBehaviour2(user, target string, act uint32) (bool, error) {
	db, err := nosql.GetBehaviourByTarget2(user, target, act)
	if err != nil && !strings.Contains(err.Error(), "no documents in result") {
		return false, err
	}
	if db != nil {
		return true, nil
	}
	return false, nil
}

func (mine *cacheContext) GetBehaviourCountByUser(user string) int64 {
	num, _ := nosql.GetBehaviourCountByUser(user)
	return num
}

func (mine *cacheContext) GetBehaviourCount(target string, act ActionType) int64 {
	num, _ := nosql.GetBehaviourCountByAction(target, uint8(act))
	return num
}

// 获取用户浏览历史数据
func (mine *cacheContext) GetBehaviourHistories(user string, kind TargetType) []*nosql.Behaviour {
	list, err := nosql.GetBehaviourByType(user, uint8(kind), 20)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return list
}

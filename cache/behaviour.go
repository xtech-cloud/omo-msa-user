package cache

import (
	"github.com/micro/go-micro/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.user/proxy/nosql"
	"sort"
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

func (mine *cacheContext) createBehaviour(user, target, scene, operator string, kind TargetType, act ActionType) error {
	db := new(nosql.Behaviour)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetBehaviourNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.UpdatedTime = time.Now()
	db.User = user
	db.Target = target
	db.Type = uint8(kind)
	db.Action = uint8(act)
	db.Scene = scene

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

func (mine *cacheContext) AddBehaviour(user, target, scene, operator string, kind TargetType, act ActionType) error {
	had, err := mine.HadBehaviour(user, target)
	if err != nil {
		return err
	}
	if had {
		return nil
	}
	return mine.createBehaviour(user, target, scene, operator, kind, act)
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

func (mine *cacheContext) GetBehavioursLatestByScene(scene string, tp, num uint32) []*nosql.Behaviour {
	list, err := nosql.GetBehavioursByScene(scene, tp, int64(num))
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return list
}

func (mine *cacheContext) GetBehaviourByTarget(user, target string) (*nosql.Behaviour, error) {
	db, err := nosql.GetBehaviourByTarget(user, target)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (mine *cacheContext) GetBehavioursByTarget(target string) ([]*nosql.Behaviour, error) {
	dbs, err := nosql.GetBehavioursByTarget(target)
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func (mine *cacheContext) GetTopBehavioursBy(users, targets []string, num uint32) []*nosql.Behaviour {
	all := make([]*nosql.Behaviour, 0, 500)
	for _, user := range users {
		for _, target := range targets {
			db, _ := nosql.GetBehaviourByTarget(user, target)
			if db != nil {
				all = append(all, db)
			}
		}
	}
	if uint32(len(all)) > num {
		sort.Slice(all, func(i, j int) bool {
			return all[i].CreatedTime.Unix() > all[j].CreatedTime.Unix()
		})
		_, _, list := CheckPage(1, num, all)
		return list
	}

	return all
}

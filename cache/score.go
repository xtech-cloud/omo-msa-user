package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.user/proxy"
	"omo.msa.user/proxy/nosql"
	"time"
)

const (
	ScoreUser ScoreType = 1
	ScoreDate ScoreType = 2
)

type ScoreType uint8

type ScoreInfo struct {
	ID          uint64 `json:"-"`
	CreatedTime time.Time
	UpdatedTime time.Time
	UID         string `json:"uid"`
	Creator     string
	Operator    string
	Name        string
	Type        uint32
	Stamp       int64

	Total  uint64
	Scene  string
	Entity string

	List []proxy.ScorePair
}

func (mine *cacheContext) CreateScore(scene, entity, operator string, count, kind uint32, stamp int64, tp ScoreType) (*ScoreInfo, error) {
	db := new(nosql.Score)
	db.Scene = scene
	db.Entity = entity
	db.ID = nosql.GetScoreNextID()
	db.UID = primitive.NewObjectID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Name = ""
	db.Stamp = stamp
	db.Type = uint8(tp)
	db.Count = uint64(count)
	db.List = make([]proxy.ScorePair, 0, 1)
	if kind > 0 {
		db.List = append(db.List, proxy.ScorePair{
			Type: kind, Count: count,
		})
	}

	err := nosql.CreateScore(db)
	if err != nil {
		return nil, err
	}
	info := new(ScoreInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) UpdateScoreByDate(scene, operator string, num uint32) error {
	now := time.Now()
	dt := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	stamp := dt.UTC().Unix()
	info := mine.GetScoreInfoByDate(scene, stamp)
	var er error
	if info == nil {
		info, er = mine.CreateScore(scene, "", operator, num, 0, stamp, ScoreDate)
		if er != nil {
			return er
		}
	}
	return info.UpdateCount(0, num, operator)
}

func (mine *cacheContext) GetScoreInfo(scene, entity string) (*ScoreInfo, error) {
	db, err := nosql.GetScoreBySceneEntity(scene, entity)
	if err != nil {
		return nil, err
	}
	info := new(ScoreInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetScores(entity string) ([]*ScoreInfo, error) {
	dbs, err := nosql.GetScoresByEntity(entity)
	if err != nil {
		return nil, err
	}
	list := make([]*ScoreInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(ScoreInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list, nil
}

func (mine *cacheContext) GetScoreInfoByDate(scene string, date int64) *ScoreInfo {
	db, err := nosql.GetScoreBySceneDate(scene, date)
	if err != nil {
		return nil
	}
	info := new(ScoreInfo)
	info.initInfo(db)
	return info
}

func (mine *cacheContext) RemoveScore(uid, operator string) error {
	err := nosql.RemoveScore(uid, operator)
	if err != nil {
		return err
	}
	return nil
}

func (mine *cacheContext) GetScoreInfoBy(uid string) (*ScoreInfo, error) {
	db, err := nosql.GetScore(uid)
	if err != nil {
		return nil, err
	}
	info := new(ScoreInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetScoresByHalfYear(scene string) []*ScoreInfo {
	now := time.Now()
	list := make([]*ScoreInfo, 0, 6)
	index := 1
	if now.Month() > 6 {
		index = 7
	}
	for i := 0; i < 6; i += 1 {
		mon := time.Month(index + i)
		from := time.Date(now.Year(), mon, 1, 0, 0, 0, 0, time.UTC).Unix()
		to := time.Date(now.Year(), mon, 30, 0, 0, 0, 0, time.UTC).Unix()
		var num uint64
		dbs, _ := nosql.GetScoreBySceneDur(scene, from, to)
		for _, db := range dbs {
			num += db.Count
		}
		list = append(list, &ScoreInfo{Scene: scene, Type: uint32(mon), Total: num})
	}
	return list
}

func (mine *cacheContext) GetScoresByWeek(scene string) []*ScoreInfo {
	now := time.Now()
	list := make([]*ScoreInfo, 0, 7)
	weekDay := int(now.Weekday())
	for i := 0; i < 7; i += 1 {
		if i <= weekDay {
			begin := now.AddDate(0, 0, -(weekDay - i))
			day := time.Date(begin.Year(), begin.Month(), begin.Day(), 0, 0, 0, 0, time.UTC)
			stamp := day.Unix()
			db, _ := nosql.GetScoreBySceneDate(scene, stamp)
			var num uint64 = 0
			if db != nil {
				num = db.Count
			}
			list = append(list, &ScoreInfo{Scene: scene, Type: uint32(i), Total: num})
		} else {
			list = append(list, &ScoreInfo{Scene: scene, Type: uint32(i), Total: 0})
		}
	}
	return list
}

func (mine *cacheContext) GetUserScoreTops(scene string, top uint32) []*ScoreInfo {
	if top < 1 {
		top = 20
	}
	var dbs []*nosql.Score
	var err error
	if scene == "" {
		dbs, err = nosql.GetScoresByTop(int64(top), uint8(ScoreUser))
	} else {
		dbs, err = nosql.GetScoresBySceneTop(scene, int64(top), uint8(ScoreUser))
	}
	if err != nil {
		return nil
	}
	list := make([]*ScoreInfo, 0, len(dbs))
	for _, db := range dbs {
		tmp := new(ScoreInfo)
		tmp.initInfo(db)
		list = append(list, tmp)
	}
	return list
}

func (mine *cacheContext) GetDateScoreTops(scene string, top uint32) []*ScoreInfo {
	if top < 1 {
		return nil
	}
	var dbs []*nosql.Score
	var err error
	if scene == "" {
		dbs, err = nosql.GetScoresByTop(int64(top), uint8(ScoreDate))
	} else {
		dbs, err = nosql.GetScoresBySceneTop(scene, int64(top), uint8(ScoreDate))
	}
	if err != nil {
		return nil
	}
	list := make([]*ScoreInfo, 0, len(dbs))
	for _, db := range dbs {
		tmp := new(ScoreInfo)
		tmp.initInfo(db)
		list = append(list, tmp)
	}
	return list
}

func (mine *ScoreInfo) initInfo(db *nosql.Score) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreatedTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Type = uint32(db.Type)
	mine.Entity = db.Entity
	mine.Scene = db.Scene
	mine.UpdatedTime = db.UpdatedTime
	mine.Total = db.Count
	mine.Stamp = db.Stamp
	mine.List = make([]proxy.ScorePair, 0, len(db.List))
	for _, pair := range mine.List {
		mine.List = append(mine.List, proxy.ScorePair{Type: pair.Type, Count: pair.Count})
	}
}

func (mine *ScoreInfo) UpdateCount(tp, num uint32, operator string) error {
	count := mine.Total + uint64(num)
	if tp > 0 {
		arr := make([]proxy.ScorePair, 0, len(mine.List))
		arr = append(arr, mine.List...)
		had := false

		for _, pair := range arr {
			if pair.Type == tp {
				pair.Count += num
				had = true
				break
			}
		}
		if !had {
			arr = append(arr, proxy.ScorePair{
				Type:  tp,
				Count: num,
			})
		}
		err := nosql.UpdateScorePairs(mine.UID, operator, count, arr)
		if err != nil {
			return err
		}
		mine.List = arr
	} else {
		count = mine.Total + uint64(num)
		err := nosql.UpdateScoreCount(mine.UID, operator, count)
		if err != nil {
			return err
		}
	}
	mine.Total = count
	mine.Operator = operator

	return nil
}

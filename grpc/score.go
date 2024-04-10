package grpc

import (
	"context"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
	"omo.msa.user/cache"
	"strconv"
)

type ScoreService struct{}

func switchScore(info *cache.ScoreInfo) *pb.ScoreInfo {
	tmp := new(pb.ScoreInfo)
	tmp.Uid = info.UID
	tmp.Updated = info.UpdatedTime.Unix()
	tmp.Created = info.CreatedTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator

	tmp.Scene = info.Scene
	tmp.Name = info.Name
	tmp.Entity = info.Entity
	tmp.Total = info.Total
	tmp.Type = info.Type
	tmp.List = make([]*pb.ScorePair, 0, len(info.List))
	for _, pair := range info.List {
		tmp.List = append(tmp.List, &pb.ScorePair{Type: pair.Type, Count: pair.Count})
	}
	return tmp
}

func (mine *ScoreService) AddOne(ctx context.Context, in *pb.ReqScoreAdd, out *pb.ReplyScoreInfo) error {
	path := "score.addOne"
	inLog(path, in)

	if in.Scene == "" {
		in.Scene = "system"
	}
	if in.Entity == "" {
		out.Status = outError(path, "the entity is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetScoreInfo(in.Scene, in.Entity)
	if info == nil {
		info, err = cache.Context().CreateScore(in.Scene, in.Entity, in.Operator, in.Count, in.Type, 0, cache.ScoreUser)
	} else {
		err = info.UpdateCount(in.Type, in.Count, in.Operator)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	_ = cache.Context().UpdateScoreByDate(in.Scene, in.Operator, in.Count)
	out.Info = switchScore(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScoreService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyScoreInfo) error {
	path := "score.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the asset is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	score, err := cache.Context().GetScoreInfoBy(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchScore(score)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScoreService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "score.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the score uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	err := cache.Context().RemoveScore(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScoreService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "score.getStatistic"
	inLog(path, in)
	if in.Key == "entity" {
		arr, _ := cache.Context().GetScores(in.Value)
		for _, info := range arr {
			out.Count += info.Total
		}
	} else if in.Key == "score" {
		tp, _ := strconv.ParseInt(in.Value, 10, 32)
		out.Count = uint64(cache.Context().GetScoreCountByType(in.Uid, uint32(tp)))
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScoreService) UpdateByFilter(ctx context.Context, in *pb.ReqUpdateFilter, out *pb.ReplyInfo) error {
	path := "score.updateByFilter"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *ScoreService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyScoreList) error {
	path := "score.getByFilter"
	inLog(path, in)
	var err error
	var list []*cache.ScoreInfo
	if in.Key == "top_user" || in.Key == "top" {
		top, _ := strconv.ParseInt(in.Value, 10, 32)
		list = cache.Context().GetUserScoreTops(in.Owner, uint32(top))
	} else if in.Key == "year_half" {
		list = cache.Context().GetScoresByHalfYear(in.Owner)
	} else if in.Key == "week" {
		//num, _ := strconv.ParseInt(in.Value, 10, 32)
		list = cache.Context().GetScoresByWeek(in.Owner)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ScoreInfo, 0, len(list))
	for _, info := range list {
		tmp := switchScore(info)
		out.List = append(out.List, tmp)
	}
	out.Status = outLog(path, out)
	return nil
}

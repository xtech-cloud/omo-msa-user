package grpc

import (
	"context"
	"fmt"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
	"omo.msa.user/cache"
	"omo.msa.user/proxy/nosql"
	"strconv"
)

type BehaviourService struct{}

func switchBehaviour(info *nosql.Behaviour) *pb.BehaviourInfo {
	tmp := &pb.BehaviourInfo{
		Uid:     info.UID.Hex(),
		Scene:   info.Scene,
		User:    info.User,
		Target:  info.Target,
		Type:    uint32(info.Type),
		Action:  uint32(info.Action),
		Created: uint64(info.CreatedTime.Unix()),
		Creator: info.Creator,
		Updated: uint64(info.UpdatedTime.Unix()),
	}
	return tmp
}

func (mine *BehaviourService) AddOne(ctx context.Context, in *pb.ReqBehaviourAdd, out *pb.ReplyInfo) error {
	path := "behaviour.addOne"
	inLog(path, in)
	var err error
	msg, _ := cache.Context().GetMessageByQuote(in.User, in.Target)
	if msg != nil {
		err = msg.Read()
	} else {
		err = cache.Context().AddBehaviour(in.User, in.Target, in.Scene, in.Operator, cache.TargetType(in.Type), cache.ActionType(in.Action))
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *BehaviourService) HadOne(ctx context.Context, in *pb.ReqBehaviourCheck, out *pb.ReplyBehaviourCheck) error {
	path := "behaviour.hadOne"
	inLog(path, in)
	var had bool
	var err error
	if in.Action > 0 {
		had, err = cache.Context().HadBehaviour2(in.User, in.Target, in.Action)
	} else {
		had, err = cache.Context().HadBehaviour(in.User, in.Target)
	}

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	if !had {
		msg, _ := cache.Context().GetMessageByQuote(in.User, in.Target)
		if msg != nil && msg.Status == cache.MessageRead {
			had = true
		}
	}
	out.Had = had
	out.Status = outLog(path, out)
	return nil
}

func (mine *BehaviourService) GetCount(ctx context.Context, in *pb.ReqBehaviourCheck, out *pb.ReplyBehaviourCheck) error {
	path := "behaviour.getCount"
	inLog(path, in)
	num := cache.Context().GetBehaviourCountByUser(in.User)
	out.Count = uint32(num)
	out.Status = outLog(path, out)
	return nil
}

func (mine *BehaviourService) UpdateOne(ctx context.Context, in *pb.ReqBehaviourUpdate, out *pb.ReplyInfo) error {
	path := "behaviour.updateOne"
	inLog(path, in)
	err := cache.Context().UpdateBehaviour(in.User, in.Target, cache.ActionType(in.Action))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *BehaviourService) GetList(ctx context.Context, in *pb.ReqBehaviourList, out *pb.ReplyBehaviourList) error {
	path := "behaviour.getList"
	inLog(path, in)
	var list []*nosql.Behaviour
	var err error
	if len(in.User) > 1 && len(in.Target) > 1 {
		list = make([]*nosql.Behaviour, 0, 1)
		tmp, _ := cache.Context().GetBehaviourByTarget(in.User, in.Target)
		if tmp != nil {
			list = append(list, tmp)
		}
	} else if len(in.User) > 1 {
		list = cache.Context().GetBehaviourHistories(in.User, cache.TargetType(in.Type))
	} else if len(in.Target) > 1 {
		list, _ = cache.Context().GetBehavioursByTarget(in.Target)
	} else {
		out.Status = outError(path, "", pbstatus.ResultStatus_DBException)
		return nil
	}

	if err != nil {
		out.Status = outError(path, "", pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.BehaviourInfo, 0, len(list))
	for _, behaviour := range list {
		out.List = append(out.List, switchBehaviour(behaviour))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *BehaviourService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyBehaviourList) error {
	path := "behaviour.getList"
	inLog(path, in)
	var list []*nosql.Behaviour
	var err error
	if in.Key == "latest" || in.Key == "top" {
		list = cache.Context().GetTopBehavioursBy(in.Values, in.List, in.Number)
	} else if in.Key == "dynamic" {
		tp, _ := strconv.ParseInt(in.Value, 10, 32)
		list = cache.Context().GetBehavioursLatestByScene(in.Owner, uint32(tp), in.Number)
	} else {
		out.Status = outError(path, "the key not defined", pbstatus.ResultStatus_Empty)
		return nil
	}

	if err != nil {
		out.Status = outError(path, "", pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.BehaviourInfo, 0, len(list))
	for _, behaviour := range list {
		out.List = append(out.List, switchBehaviour(behaviour))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *BehaviourService) GetStatistic(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStatistic) error {
	path := "behaviour.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

package grpc

import (
	"context"
	"fmt"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
	"omo.msa.user/cache"
	"omo.msa.user/proxy/nosql"
)

type BehaviourService struct {}

func switchBehaviour(info *nosql.Behaviour) *pb.BehaviourInfo {
	tmp := &pb.BehaviourInfo{
		User : info.User,
		Target : info.Target,
		Type : uint32(info.Type),
		Action : uint32(info.Action),
		Created : uint64(info.CreatedTime.Unix()),
		Creator: info.Creator,
		Updated: uint64(info.UpdatedTime.Unix()),
	}
	return tmp
}

func (mine *BehaviourService)AddOne(ctx context.Context, in *pb.ReqBehaviourAdd, out *pb.ReplyInfo) error {
	path := "behaviour.addOne"
	inLog(path, in)
	err := cache.Context().AddBehaviour(in.User, in.Target, cache.TargetType(in.Type), cache.ActionType(in.Action))
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *BehaviourService)HadOne(ctx context.Context, in *pb.ReqBehaviourCheck, out *pb.ReplyBehaviourCheck) error {
	path := "behaviour.hadOne"
	inLog(path, in)
	had, err := cache.Context().HadBehaviour(in.User, in.Target)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Had = had
	out.Status = outLog(path, out)
	return nil
}

func (mine *BehaviourService)GetCount(ctx context.Context, in *pb.ReqBehaviourCheck, out *pb.ReplyBehaviourCheck) error {
	path := "behaviour.getCount"
	inLog(path, in)
	num := cache.Context().GetBehaviourCountByUser(in.User)
	out.Count = uint32(num)
	out.Status = outLog(path, out)
	return nil
}

func (mine *BehaviourService)UpdateOne(ctx context.Context, in *pb.ReqBehaviourUpdate, out *pb.ReplyInfo) error {
	path := "behaviour.updateOne"
	inLog(path, in)
	err := cache.Context().UpdateBehaviour(in.User, in.Target, cache.ActionType(in.Action))
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *BehaviourService)GetList(ctx context.Context, in *pb.ReqBehaviourList, out *pb.ReplyBehaviourList) error {
	path := "behaviour.getList"
	inLog(path, in)
	var list []*nosql.Behaviour
	var err error
	if len(in.User) > 1 && len(in.Target) > 1 {

	}else if len(in.User) > 1 {
		list = cache.Context().GetBehaviourHistories(in.User, cache.TargetType(in.Type))
	}else if len(in.Target) > 1 {

	}else {
		out.Status = outError(path,"", pbstatus.ResultStatus_DBException)
		return nil
	}

	if err != nil {
		out.Status = outError(path,"", pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.BehaviourInfo, 0, len(list))
	for _, behaviour := range list {
		out.List = append(out.List, switchBehaviour(behaviour))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

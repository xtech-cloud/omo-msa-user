package grpc

import (
	"context"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
	"omo.msa.user/cache"
)

type WechatService struct{}

func switchWechat(info *cache.WechatInfo) *pb.WechatInfo {
	tmp := &pb.WechatInfo{
		Uid:      info.UID,
		Id:       info.ID,
		Created:  uint64(info.CreateTime.Unix()),
		Updated:  uint64(info.UpdateTime.Unix()),
		Operator: info.Operator,
		Creator:  info.Creator,
		Open:     info.OpenID,
		Union:    info.UnionID,
		Portrait: info.Portrait,
	}
	return tmp
}

func (mine *WechatService) AddOne(ctx context.Context, in *pb.ReqWechatAdd, out *pb.ReplyWechatInfo) error {
	path := "wechat.addOne"
	inLog(path, in)
	info, err := cache.Context().CreateWechat(in.Name, in.Open, in.Union, in.Portrait, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchWechat(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *WechatService) GetOne(ctx context.Context, in *pb.ReqWechatBy, out *pb.ReplyWechatInfo) error {
	path := "wechat.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	var info *cache.WechatInfo
	if in.Type == pb.WechatType_Default {
		info = cache.Context().GetWechat(in.Uid)
	} else if in.Type == pb.WechatType_OpenID {
		info = cache.Context().GetWechatByOpen(in.Uid)
	} else if in.Type == pb.WechatType_Union {

	}

	if info == nil {
		out.Status = outError(path, "the wechat not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Info = switchWechat(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *WechatService) UpdateBase(ctx context.Context, in *pb.ReqWechatUpdate, out *pb.ReplyWechatInfo) error {
	path := "wechat.updateOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetWechat(in.Uid)
	if info == nil {
		out.Status = outError(path, "the wechat not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateBase(in.Name, in.Open, in.Union, in.Portrait, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchWechat(info)
	out.Status = outLog(path, out)
	return nil
}

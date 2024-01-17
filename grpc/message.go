package grpc

import (
	"context"
	"fmt"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
	"omo.msa.user/cache"
)

type MessageService struct{}

func switchMessage(info *cache.MessageInfo) *pb.MessageInfo {
	tmp := &pb.MessageInfo{
		Uid:     info.UID,
		User:    info.User,
		Quote:   info.Quote,
		Type:    uint32(info.Type),
		Status:  uint32(info.Status),
		Stamp:   info.Stamp,
		Owner:   info.Owner,
		Targets: info.Targets,
		Created: uint64(info.CreateTime.Unix()),
		Creator: info.Creator,
		Updated: uint64(info.UpdateTime.Unix()),
	}
	return tmp
}

func (mine *MessageService) AddOne(ctx context.Context, in *pb.ReqMessageAdd, out *pb.ReplyInfo) error {
	path := "message.addOne"
	inLog(path, in)
	if in.User == "" || in.Quote == "" {
		out.Status = outError(path, "the user or quote is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, _ := cache.Context().GetMessageByQuote(in.User, in.Quote)
	if info == nil {
		err := cache.Context().CreateMessage(in.Owner, in.User, in.Quote, in.Operator, in.Type, in.Stamp, in.Targets)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	} else {
		err := info.UpdateTargets(in.Targets)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *MessageService) GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyMessageList) error {
	path := "message.getList"
	inLog(path, in)
	var list []*cache.MessageInfo
	//var err error
	if in.Key == "user" {
		list = cache.Context().GetMessagesByUser(in.Value)
	} else if in.Key == "quote" {
		msg, _ := cache.Context().GetMessageByQuote(in.Param, in.Value)
		if msg != nil {
			list = make([]*cache.MessageInfo, 0, 1)
			list = append(list, msg)
		}
	} else {
		out.Status = outError(path, "", pbstatus.ResultStatus_DBException)
		return nil
	}

	//if err != nil {
	//	out.Status = outError(path, "", pbstatus.ResultStatus_DBException)
	//	return nil
	//}
	t, p, arr := cache.CheckPage(in.Page, in.Number, list)
	out.List = make([]*pb.MessageInfo, 0, len(list))
	for _, message := range arr {
		_ = message.Awake()
		out.List = append(out.List, switchMessage(message))
	}
	out.Total = t
	out.Pages = p
	out.Status = outLog(path, fmt.Sprintf("the now length = %d and total = %d", len(out.List), t))
	return nil
}

func (mine *MessageService) GetStatistic(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStatistic) error {
	path := "message.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *MessageService) UpdateByFilter(ctx context.Context, in *pb.ReqUpdateFilter, out *pb.ReplyInfo) error {
	path := "message.updateByFilter"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

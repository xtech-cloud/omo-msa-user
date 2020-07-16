package grpc

import (
	"context"
	"errors"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
	"omo.msa.user/cache"
)

type UserService struct {}

func switchUser(info *cache.UserInfo) *pb.UserInfo {
	tmp := new(pb.UserInfo)
	tmp.Job = info.Datum.Job
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Type = pb.UserType(info.Type)
	tmp.Account = info.Account
	tmp.Sex = pb.UserSex(info.Datum.Sex)
	tmp.Phone = info.Datum.Phone
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.RealName = info.Datum.RealName
	return tmp
}

func (mine *UserService)AddOne(ctx context.Context, in *pb.ReqUserAdd, out *pb.ReplyUserOne) error {
	if len(in.Account) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the account is empty")
	}
	info := new(cache.DatumInfo)
	info.RealName = in.Name
	info.Phone = in.Phone
	info.Sex = uint8(in.Sex)
	user,err := cache.CreateUser(in.Account, in.NickName, in.Remark, uint8(in.Type), info)
	if err == nil {
		out.Info = switchUser(user)
	}else{
		out.ErrorCode = pb.ResultStatus_DBException
	}

	return err
}

func (mine *UserService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyUserOne) error {
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetUser(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the user not found")
	}
	out.Info = switchUser(info)
	return nil
}

func (mine *UserService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	err := cache.RemoveUser(in.Uid, in.Operator)
	if err != nil {
		out.ErrorCode = pb.ResultStatus_DBException
	}
	return err
}

func (mine *UserService)GetList(ctx context.Context, in *pb.ReqUserList, out *pb.ReplyUserList) error {
	out.List = make([]*pb.UserInfo, 0, len(in.List))
	for _, value := range in.List {
		info := cache.GetUser(value)
		if info != nil {
			out.List = append(out.List, switchUser(info))
		}
	}
	return nil
}

func (mine *UserService)GetByPage(ctx context.Context, in *pb.RequestPage, out *pb.ReplyUserList) error {
	out.List = make([]*pb.UserInfo, 0, in.Number)
	users := cache.AllUsers()
	total := uint32(len(users))
	out.PageMax = total / in.Number + 1
	var i uint32 = 0
	for ;i < total;i += 1{
		t := i / in.Number + 1
		if t == in.Page {
			out.List = append(out.List, switchUser(users[i]))
		}
	}
	out.PageNow = in.Page
	out.Total = uint64(total)
	return nil
}

func (mine *UserService) GetByAccount (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyUserOne) error {
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetUserByAccount(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the user not found")
	}
	out.Info = switchUser(info)
	return nil
}

func (mine *UserService) UpdateBase (ctx context.Context, in *pb.ReqUserUpdate, out *pb.ReplyInfo) error {
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetUser(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the user not found")
	}
	info.UpdateBase(in.NickName, in.Name, in.Phone, in.Remark, in.Job, in.Operator, uint8(in.Sex))
	out.Uid = in.Uid
	return nil
}


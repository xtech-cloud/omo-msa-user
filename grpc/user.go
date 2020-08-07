package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/micro/go-micro/v2/logger"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
	"omo.msa.user/cache"
)

type UserService struct {}

func switchUser(info *cache.UserInfo) *pb.UserInfo {
	tmp := &pb.UserInfo{
		Uid : info.UID,
		Id : info.ID,
		Type : pb.UserType(info.Type),
		Account : info.Account,
		Sex : pb.UserSex(info.Datum.Sex),
		Phone : info.Datum.Phone,
		Name : info.Name,
		Remark : info.Remark,
		Created : info.CreateTime.Unix(),
		Updated : info.UpdateTime.Unix(),
		RealName : info.Datum.RealName,
	}
	return tmp
}

func inLog(name, data interface{})  {
	bytes, _ := json.Marshal(data)
	msg := ByteString(bytes)
	logger.Infof("[in.%s]:data = %s", name, msg)
}

func ByteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

func (mine *UserService)AddOne(ctx context.Context, in *pb.ReqUserAdd, out *pb.ReplyUserOne) error {
	inLog("user.add", in)
	if len(in.Account) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the account is empty")
	}
	tmp := cache.GetUserByAccount(in.Account)
	if tmp != nil {
		out.Info = switchUser(tmp)
		return nil
	}
	tmp1 := cache.GetUserByPhone(in.Phone)
	if tmp1 != nil {
		out.Info = switchUser(tmp1)
		return nil
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
	inLog("user.getOne", in)
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
	inLog("user.remove", in)
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
	inLog("user.account", in)
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

func (mine *UserService) GetByPhone (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyUserOne) error {
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

func (mine *UserService) UpdateBase (ctx context.Context, in *pb.ReqUserUpdate, out *pb.ReplyUserOne) error {
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetUser(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the user not found")
	}
	err := info.UpdateBase(in.NickName, in.Name, in.Phone, in.Remark, in.Job, in.Operator, uint8(in.Sex))
	if err == nil {
		out.Info = switchUser(info)
	}

	return err
}

func (mine *UserService) UpdatePasswords (ctx context.Context, in *pb.ReqUserPasswords, out *pb.ReplyInfo) error {
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetUser(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the user not found")
	}
	err := info.UpdatePasswords(in.Passwords, in.Operator)
	return err
}

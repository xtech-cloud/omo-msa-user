package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
	"omo.msa.user/cache"
)

type AccountService struct {}

func switchAccount(info *cache.AccountInfo) *pb.AccountInfo {
	tmp := &pb.AccountInfo{
		Uid : info.UID,
		Name : info.Name,
		Passwords : info.Passwords,
		Created : info.CreateTime.Unix(),
		Updated : info.UpdateTime.Unix(),
		Creator: info.Creator,
		Operator: info.Operator,
	}
	return tmp
}

func (mine *AccountService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAccount) error {
	path := "account.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetAccount(in.Uid)
	if info == nil {
		out.Status = outError(path,"the account not found", pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchAccount(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AccountService)SignIn(ctx context.Context, in *pb.ReqSignIn, out *pb.ReplyInfo) error {
	path := "account.signIn"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the account name is empty", pb.ResultCode_Empty)
		return nil
	}
	user,err := cache.Context().SignIn(in.Name, in.Psw)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Uid = user
	out.Status = outLog(path, out)
	return nil
}

func (mine *AccountService)SetPasswords(ctx context.Context, in *pb.ReqSetPasswords, out *pb.ReplyInfo) error {
	path := "account.setPasswords"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetAccount(in.Uid)
	if info == nil {
		out.Status = outError(path,"the account not found", pb.ResultCode_DBException)
		return nil
	}
	err := info.UpdatePasswords(in.Psw, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

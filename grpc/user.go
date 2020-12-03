package grpc

import (
	"context"
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
		Sex : pb.UserSex(info.Sex),
		Phone : info.Phone,
		Name : info.Name,
		Remark : info.Remark,
		Created : info.CreateTime.Unix(),
		Updated : info.UpdateTime.Unix(),
		Operator: info.Operator,
		Creator: info.Creator,
		Nick : info.NickName,
		Portrait: info.Portrait,
		Entity: info.Entity,
	}
	return tmp
}

func (mine *UserService)AddOne(ctx context.Context, in *pb.ReqUserAdd, out *pb.ReplyUserOne) error {
	path := "user.add"
	inLog(path, in)
	var err error
	var account *cache.AccountInfo
	if len(in.Account) > 1 {
		account = cache.Context().GetAccount(in.Account)
		if account == nil {
			out.Status = outError(path,"the account not found ", pb.ResultCode_NotExisted)
			return nil
		}
	}else{
		if in.Type == pb.UserType_SuperRoot {
			account, err = cache.Context().CreateAccount(in.Name, in.Passwords, in.Operator)
		}else if in.Type == pb.UserType_EnterpriseAdmin || in.Type == pb.UserType_EnterpriseCommon {
			account, err = cache.Context().CreateAccount(in.Phone, in.Passwords, in.Operator)
		}else{
			name := in.Phone
			if name == ""{
				name = in.Name
			}
			account, err = cache.Context().CreateAccount(name, in.Passwords, in.Operator)
		}
		if err != nil {
			out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
			return nil
		}
	}

	user,err1 := account.CreateUser(in.Name, in.Remark, in.Nick, in.Phone, in.Entity,in.Portrait, in.Operator, uint8(in.Type), uint8(in.Sex))
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pb.ResultCode_DBException)
		return nil
	}
	if in.Sns != nil && len(in.Sns.Uid) > 2 {
		er := user.AppendSNS(in.Sns.Uid, in.Sns.Name, uint8(in.Sns.Type))
		if er != nil {
			out.Status = outError(path,er.Error(), pb.ResultCode_DBException)
			return nil
		}
	}
	out.Info = switchUser(user)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyUserOne) error {
	path := "user.get"
	inLog(path, in)
	var info *cache.UserInfo
	if len(in.Uid) > 0 {
		info = cache.Context().GetUser(in.Uid)
	}else if len(in.Entity) > 0 {
		info = cache.Context().GetUserByEntity(in.Entity)
	}
	if info == nil {
		out.Status = outError(path,"the user not found ", pb.ResultCode_NotExisted)
		return nil
	}

	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "user.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultCode_Empty)
		return nil
	}
	account := cache.Context().GetAccountByUser(in.Uid)
	if account == nil {
		out.Status = outError(path,"the account not found ", pb.ResultCode_NotExisted)
		return nil
	}
	err := account.RemoveUser(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return err
}

func (mine *UserService)GetList(ctx context.Context, in *pb.ReqUserList, out *pb.ReplyUserList) error {
	path := "user.list"
	inLog(path, in)
	out.List = make([]*pb.UserInfo, 0, len(in.List))
	for _, value := range in.List {
		info := cache.Context().GetUser(value)
		if info != nil {
			out.List = append(out.List, switchUser(info))
		}
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService)GetByPage(ctx context.Context, in *pb.RequestPage, out *pb.ReplyUserList) error {
	path := "user.getByPage"
	inLog(path, in)
	out.List = make([]*pb.UserInfo, 0, in.Number)
	users := cache.Context().AllUsers()
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
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) GetBySNS (ctx context.Context, in *pb.ReqUserBy, out *pb.ReplyUserOne) error {
	path := "user.getBySNS"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the sns uid is empty ", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetUserBySNS(in.Uid, uint8(in.Type))
	if info == nil {
		out.Status = outError(path,"the user not found ", pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) GetByPhone (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyUserOne) error {
	path := "user.getByPhone"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the phone is empty ", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetUserByPhone(in.Uid)
	if info == nil {
		out.Status = outError(path,"the user not found ", pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) GetByID (ctx context.Context, in *pb.RequestIDInfo, out *pb.ReplyUserOne) error {
	path := "user.getByID"
	inLog(path, in)
	if in.Id < 1 {
		out.Status = outError(path,"the user id is empty ", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetUserByID(in.Id)
	if info == nil {
		out.Status = outError(path,"the user not found ", pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) GetByEntity (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyUserOne) error {
	path := "user.getByPhone"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the entity is empty ", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetUserByEntity(in.Uid)
	if info == nil {
		out.Status = outError(path,"the user not found ", pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) UpdateBase (ctx context.Context, in *pb.ReqUserUpdate, out *pb.ReplyUserOne) error {
	path := "user.update"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetUser(in.Uid)
	if info == nil {
		out.Status = outError(path,"the user not found ", pb.ResultCode_NotExisted)
		return nil
	}
	err := info.UpdateBase(in.Name, in.NickName, in.Remark, in.Portrait, in.Operator, uint8(in.Sex))
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return err
}

func (mine *UserService) UpdateEntity (ctx context.Context, in *pb.ReqUserEntity, out *pb.ReplyUserOne) error {
	path := "user.updateEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetUser(in.Uid)
	if info == nil {
		out.Status = outError(path,"the user not found ", pb.ResultCode_NotExisted)
		return nil
	}
	err := info.UpdateEntity(in.Entity, "")
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return err
}

func (mine *UserService) UpdateSNS (ctx context.Context, in *pb.ReqUserSNS, out *pb.ReplyUserOne) error {
	path := "user.updateSNS"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetUser(in.User)
	if info == nil {
		out.Status = outError(path,"the user not found ", pb.ResultCode_NotExisted)
		return nil
	}
	var err error
	if in.Add {
		err = info.AppendSNS(in.Uid, in.Name, uint8(in.Type))
	}else{
		err = info.SubtractSNS(in.Uid)
	}
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return err
}

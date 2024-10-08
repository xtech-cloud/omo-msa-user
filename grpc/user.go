package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/micro/go-micro/v2/logger"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
	"omo.msa.user/cache"
	"strings"
)

type UserService struct{}

func switchUser(info *cache.UserInfo) *pb.UserInfo {
	tmp := &pb.UserInfo{
		Uid:      info.UID,
		Id:       info.ID,
		Type:     pb.UserType(info.Type),
		Account:  info.Account,
		Sex:      pb.UserSex(info.Sex),
		Phone:    info.Phone,
		Name:     info.Name,
		Remark:   info.Remark,
		Created:  info.CreateTime.Unix(),
		Updated:  info.UpdateTime.Unix(),
		Operator: info.Operator,
		Creator:  info.Creator,
		Nick:     info.NickName,
		Portrait: info.Portrait,
		Entity:   info.Entity,
		Follows:  info.Follows,
		Relates:  info.Relates,
		Status:   uint32(info.Status),
		Tags:     info.Tags,
	}
	tmp.Shown = &pb.ShownInfo{Name: info.Shown.Name, Cover: info.Shown.Cover}
	//if info.Shown != nil {
	//	tmp.Shown = &pb.ShownInfo{Name: info.Shown.Name, Cover: info.Shown.Cover}
	//} else {
	//	//tmp.Shown = &pb.ShownInfo{Name: info.NickName, Cover: info.Portrait}
	//	tmp.Shown = &pb.ShownInfo{Name: "", Cover: ""}
	//}
	tmp.Sns = make([]*pb.SNSInfo, 0, len(info.SNS))
	for _, sn := range info.SNS {
		tmp.Sns = append(tmp.Sns, &pb.SNSInfo{Uid: sn.UID, Type: uint32(sn.Type), Name: sn.Name})
	}
	return tmp
}

func (mine *UserService) AddOne(ctx context.Context, in *pb.ReqUserAdd, out *pb.ReplyUserOne) error {
	path := "user.addOne"
	inLog(path, in)
	var err error
	var account *cache.AccountInfo
	in.Name = strings.TrimSpace(in.Name)
	in.Nick = strings.TrimSpace(in.Nick)
	if len(in.Account) > 1 {
		account = cache.Context().GetAccount(in.Account)
		if account == nil {
			out.Status = outError(path, "the account not found ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
	} else {
		psw := cache.CryptPsw(in.Passwords)
		//if in.Type == pb.UserType_SuperRoot {
		//	account, err = cache.Context().CreateAccount(in.Name, psw, in.Operator)
		//} else if in.Type == pb.UserType_EnterpriseAdmin || in.Type == pb.UserType_EnterpriseCommon {
		//	account, err = cache.Context().CreateAccount(in.Phone, psw, in.Operator)
		//} else {
		//
		//}
		name := in.Phone
		if name == "" {
			if in.Sns != nil && len(in.Sns.Uid) > 0 {
				name = in.Sns.Uid
			} else {
				name = in.Name
			}
		}
		account, err = cache.Context().CreateAccount(name, psw, in.Operator)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}

	user, err1 := account.CreateUser(in)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	if len(user.Phone) < 1 && len(in.Phone) > 1 {
		_ = user.UpdatePhone(in.Phone, in.Operator)
	}
	if uint8(in.Type) < user.Type && in.Type > 0 {
		_ = user.UpdateType(uint8(in.Type))
	}
	if in.Sns != nil && len(in.Sns.Uid) > 2 {
		er := user.AppendSNS(in.Sns.Uid, in.Sns.Name, uint8(in.Sns.Type))
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}
	out.Info = switchUser(user)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyUserOne) error {
	path := "user.getOne"
	inLog(path, in)
	var info *cache.UserInfo
	if len(in.Uid) > 0 {
		info = cache.Context().GetUser(in.Uid)
	} else if len(in.Entity) > 0 {
		info = cache.Context().GetUserByEntity(in.Entity)
	} else if len(in.Operator) > 0 {
		info = cache.Context().GetUserByName(in.Operator)
	}
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "user.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	//out.Status = outError(path,"the user cannot remove ", pbstatus.ResultStatus_DBException)
	//return nil
	err := cache.Context().RemoveUser(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) GetList(ctx context.Context, in *pb.ReqUserList, out *pb.ReplyUserList) error {
	path := "user.list"
	inLog(path, in)
	out.List = make([]*pb.UserInfo, 0, len(in.List))
	for _, value := range in.List {
		info := cache.Context().GetUser(value)
		if info != nil {
			out.List = append(out.List, switchUser(info))
		} else {
			logger.Warn("not found the user(" + value + ")")
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *UserService) GetByPage(ctx context.Context, in *pb.RequestPage, out *pb.ReplyUserList) error {
	path := "user.getByPage"
	inLog(path, in)
	out.List = make([]*pb.UserInfo, 0, in.Number)
	var list []*cache.UserInfo
	var total uint32
	var pages uint32
	if in.Key == "" {
		total, pages, list = cache.Context().GetUsersByPage(in.Page, in.Number)
	} else if in.Key == "scene" {
		total, pages, list = cache.Context().GetUsersByPageScene(in.Value, in.Page, in.Number)
	} else if in.Key == "latest" {
		total, pages, list = cache.Context().GetUsersByLatest(in.Value, in.Page, in.Number)
	}
	for _, info := range list {
		out.List = append(out.List, switchUser(info))
	}
	out.PageNow = in.Page
	out.PageMax = pages
	out.Total = uint64(total)
	out.Status = outLog(path, fmt.Sprintf("the total = %d, pages = %d and length = %d", total, pages, len(out.List)))
	return nil
}

func (mine *UserService) GetByKey(ctx context.Context, in *pb.ReqUserSearch, out *pb.ReplyUserList) error {
	path := "user.getByKey"
	inLog(path, in)

	users := cache.Context().SearchUsers(in.Type, in.Tags)
	out.List = make([]*pb.UserInfo, 0, len(users))
	for i := 0; i < len(users); i += 1 {
		out.List = append(out.List, switchUser(users[i]))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *UserService) GetStatistic(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStatistic) error {
	path := "user.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) GetBySNS(ctx context.Context, in *pb.ReqUserBy, out *pb.ReplyUserOne) error {
	path := "user.getBySNS"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the sns uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetUserBySNS(in.Uid, uint8(in.Type))
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) GetByPhone(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyUserOne) error {
	path := "user.getByPhone"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the phone is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetUserByPhone(in.Uid)
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) GetByID(ctx context.Context, in *pb.RequestIDInfo, out *pb.ReplyUserOne) error {
	path := "user.getByID"
	inLog(path, in)
	if in.Id < 1 {
		out.Status = outError(path, "the user id is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetUserByID(in.Id)
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) GetByEntity(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyUserOne) error {
	path := "user.getByEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the entity is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetUserByEntity(in.Uid)
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) UpdateBase(ctx context.Context, in *pb.ReqUserUpdate, out *pb.ReplyUserOne) error {
	path := "user.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	in.Name = strings.TrimSpace(in.Name)
	in.NickName = strings.TrimSpace(in.NickName)
	info := cache.Context().GetUser(in.Uid)
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateBase(in.Name, in.NickName, in.Remark, in.Portrait, in.Operator, uint8(in.Sex))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) UpdateEntity(ctx context.Context, in *pb.ReqUserEntity, out *pb.ReplyUserOne) error {
	path := "user.updateEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetUser(in.Uid)
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateEntity(in.Entity, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) UpdateTags(ctx context.Context, in *pb.ReqUserTags, out *pb.ReplyUserOne) error {
	path := "user.updateTags"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetUser(in.Uid)
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateTags(in.Operator, in.Tags)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) UpdateSNS(ctx context.Context, in *pb.ReqUserSNS, out *pb.ReplyUserOne) error {
	path := "user.updateSNS"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetUser(in.User)
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Add {
		err = info.AppendSNS(in.Uid, in.Name, uint8(in.Type))
	} else {
		err = info.SubtractSNS(in.Uid)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) UpdatePhone(ctx context.Context, in *pb.ReqUserPhone, out *pb.ReplyInfo) error {
	path := "user.updatePhone"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetUser(in.Uid)
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	old := info.Phone
	if in.Phone == old {
		out.Status = outLog(path, out)
		return nil
	}
	err := info.UpdatePhone(in.Phone, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	account := cache.Context().GetAccount(info.Account)
	if account != nil && account.Name == old {
		_ = account.UpdateName(in.Phone, in.Operator)
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) UpdateFollows(ctx context.Context, in *pb.RequestList, out *pb.ReplyInfo) error {
	path := "user.updateFollows"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetUser(in.Uid)
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateFollows(in.List)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) UpdateByFilter(ctx context.Context, in *pb.ReqUpdateFilter, out *pb.ReplyInfo) error {
	path := "user.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetUser(in.Uid)
	if info == nil {
		out.Status = outError(path, "the user not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Key == "relates" {
		err = info.UpdateRelates(in.Values)
	} else if in.Key == "shown" {
		if len(in.Values) == 2 {
			err = info.UpdateShown(in.Values[0], in.Values[1])
		}
	} else {
		err = errors.New("the key not defined")
	}

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

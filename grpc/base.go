package grpc

import (
	"encoding/json"
	"github.com/micro/go-micro/v2/logger"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-user/proto/user"
)

func inLog(name, data interface{}) {
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

func outError(name, msg string, code pbstatus.ResultStatus) *pb.ReplyStatus {
	logger.Warnf("[error.%s]:code = %d, msg = %s", name, code, msg)
	tmp := &pb.ReplyStatus{
		Code: pb.ResultCode(code),
		Msg:  msg,
	}
	return tmp
}

func outLog(name, data interface{}) *pb.ReplyStatus {
	bytes, _ := json.Marshal(data)
	msg := ByteString(bytes)
	logger.Infof("[out.%s]:data = %s", name, msg)
	tmp := &pb.ReplyStatus{
		Code: 0,
		Msg:  "",
	}
	return tmp
}

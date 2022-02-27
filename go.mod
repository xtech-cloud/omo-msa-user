module omo.msa.user

go 1.15

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.26.0

require (
	github.com/labstack/gommon v0.3.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/config/source/consul/v2 v2.9.1
	github.com/micro/go-plugins/logger/logrus/v2 v2.9.1
	github.com/micro/go-plugins/registry/consul/v2 v2.9.1
	github.com/micro/go-plugins/registry/etcdv3/v2 v2.9.1
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/tidwall/gjson v1.6.0
	github.com/xtech-cloud/omo-msp-user v1.3.9
	github.com/xtech-cloud/omo-msp-status v1.0.1
	go.mongodb.org/mongo-driver v1.4.6
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

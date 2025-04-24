package grpcpayment

import (
	typesm "authentication_service/auth_service/types"
	"authentication_service/core/lib/external/grpccore"
	protoobj "authentication_service/core/proto"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type AuthServiceServiceProto struct {
	ipc *typesm.InternalProviderControl
	protoobj.UnimplementedAuthServiceServer
}

func newPaymentServiceProtoServer(ipc *typesm.InternalProviderControl) protoobj.AuthServiceServer {
	return &AuthServiceServiceProto{
		ipc: ipc,
	}
}

func StartGrpcServer(ipc *typesm.InternalProviderControl) {
	addressService := fmt.Sprintf(":%d", ipc.Config.ExposedServiceConfig.AuthService.GrpcPort)
	server, err := grpccore.CreateServerGRPC(nil, nil)
	if err != nil {
		logrus.Errorln("🔴 Failed to create gRPC server: ", err)
		return
	}

	protoobj.RegisterAuthServiceServer(server, newPaymentServiceProtoServer(ipc))

	lis, err := net.Listen("tcp", addressService)
	if err != nil {
		logrus.Errorln("🔴 Failed to listen: ", err)
		return
	}
	logrus.Infof("🟢 Server grpc listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		logrus.Errorln("🔴 Failed to serve: ", err)
		return
	}
	return
}

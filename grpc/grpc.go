package grpc

import (
	"context"
	"github.com/junminhong/member-services-center/db/redis"
	"github.com/junminhong/member-services-center/grpc/proto"
	"github.com/junminhong/member-services-center/pkg/jwt"
	"github.com/junminhong/member-services-center/pkg/logger"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var sugar = logger.Setup()
var redisClient = redis.Setup()

type Server struct{}

func (s Server) VerifyAccessToken(ctx context.Context, request *proto.TokenAuthRequest) (*proto.TokenAuthResponse, error) {
	if !jwt.VerifyAtomicToken(request.Token) {
		return nil, nil
	}
	memberID, err := redisClient.Get(context.Background(), request.Token).Result()
	if err != nil {
		sugar.Info(err.Error())
	}
	response := &proto.TokenAuthResponse{MemberID: memberID}
	return response, nil
}

func SetupServer(intiServerWg *sync.WaitGroup) {
	defer intiServerWg.Done()
	sugar.Info("starting gRPC server...")
	sugar.Info("Listening and serving HTTP on :127.0.0.1:2021")

	lis, err := net.Listen("tcp", "127.0.0.1:2021")
	if err != nil {
		sugar.Info("failed to listen: %v \n", err)
	}
	gRpcServer := grpc.NewServer()
	proto.RegisterTokenAuthServiceServer(gRpcServer, &Server{})

	if err := gRpcServer.Serve(lis); err != nil {
		sugar.Info("failed to serve: %v \n", err)
	}
}

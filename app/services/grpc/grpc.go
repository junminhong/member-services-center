package grpc

import (
	"context"
	"github.com/junminhong/member-services-center/app/proto/member"
	"github.com/junminhong/member-services-center/app/services/jwt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type Server struct{}

func (s Server) VerifyAccessToken(ctx context.Context, request *member.TokenAuthRequest) (*member.TokenAuthResponse, error) {
	response := &member.TokenAuthResponse{Response: jwt.VerifyAccessToken(request.Token)}
	return response, nil
}

func InitGRpcServer(intiServerWg *sync.WaitGroup) {
	defer intiServerWg.Done()
	log.Println("starting gRPC server...")
	log.Println("Listening and serving HTTP on :127.0.0.1:2021")

	lis, err := net.Listen("tcp", "127.0.0.1:2021")
	if err != nil {
		log.Fatalf("failed to listen: %v \n", err)
	}
	gRpcServer := grpc.NewServer()
	member.RegisterTokenAuthServiceServer(gRpcServer, &Server{})

	if err := gRpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v \n", err)
	}
}

package main

import (
	"net"

	"github.com/labstack/gommon/log"
	pb "github.com/rod6/tron/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	name = "user"
	port = ":13001"
)

type userServer struct{}

func (u *userServer) Auth(ctx context.Context, aq *pb.AuthRequest) (*pb.AuthReply, error) {
	if aq.Username == "rod" && aq.Password == "123456" {
		return &pb.AuthReply{Authed: true}, nil
	}
	return &pb.AuthReply{Authed: false}, nil
}

func main() {
	log := log.New(name)
	log.Infof("service %s is starting at port %s", name, port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServer(s, &userServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

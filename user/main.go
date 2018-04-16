package main

import (
	"net"

	"github.com/labstack/gommon/log"
	pb "github.com/rod6/tron/pb"
	. "github.com/rod6/tron/user/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	name = "user"
	port = ":13001"
)

func main() {
	Logger.Infof("service %s is starting at port %s", name, port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServer(s, &UserServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

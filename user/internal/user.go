package internal

import (
	"github.com/labstack/gommon/log"
	// . "github.com/rod6/tron/lib"
	pb "github.com/rod6/tron/pb"
	"golang.org/x/net/context"
)

var (
	Logger = log.New("user")
)

type UserServer struct{}

func (u *UserServer) Auth(c context.Context, aq *pb.AuthRequest) (*pb.AuthReply, error) {
	Logger.Infof("request id = %s", GetRid(c))

	if aq.Username == "rod" && aq.Password == "123456" {
		return &pb.AuthReply{Authed: true}, nil
	}
	return &pb.AuthReply{Authed: false}, nil
}

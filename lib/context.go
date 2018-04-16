package lib

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func GetRid(c context.Context) string {
	nilId := "00000000-0000-0000-0000-000000000000"

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nilId
	}

	rids := md["rid"]
	if len(rids) == 0 || rids[0] == "" {
		return nilId
	}

	return rids[0]
}

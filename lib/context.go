package lib

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// GetRid extract request id from context
func GetRid(c context.Context) string {
	nilID := "00000000-0000-0000-0000-000000000000"

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nilID
	}

	rids := md["rid"]
	if len(rids) == 0 || rids[0] == "" {
		return nilID
	}

	return rids[0]
}

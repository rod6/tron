package ctrl

import (
	"errors"
	"time"

	"google.golang.org/grpc"
)

var (
	services = map[string]string{
		"user": "127.0.0.1:13001",
	}

	// Connections contain the client connections to micro services
	Connections = make(map[string]*grpc.ClientConn)
)

// ConnSrv connects micro service used in gateway
func ConnSrv() error {
	for name, addr := range services {
		conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(time.Second), grpc.WithBlock())
		if err != nil {
			Logger.Errorf(`connect to '%s:%s' service failed: %v`, name, addr, err)
			return errors.New("init services connection error")
		}
		Connections[name] = conn
	}
	return nil
}

// CloseConn closes all connects to services
func CloseConn() {
	for _, conn := range Connections {
		if conn != nil {
			conn.Close()
		}
	}
}

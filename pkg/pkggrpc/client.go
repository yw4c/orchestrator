package pkggrpc

import (
	"google.golang.org/grpc"
	"orchestrator/pb"
	"sync"
)

const targetGrpc = "burnner:10000"

var pingPongCliOnce sync.Once
var pingPongCli pb.PingPongServiceClient

func GetPingPongCliInstance() pb.PingPongServiceClient {

	pingPongCliOnce.Do(func() {
		conn, err := grpc.Dial(targetGrpc, grpc.WithInsecure())
		if err != nil {
			panic(err.Error())
		}
		pingPongCli = pb.NewPingPongServiceClient(conn)
	})
	return pingPongCli
}

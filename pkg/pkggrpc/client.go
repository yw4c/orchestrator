package pkggrpc

import (
	"orchestrator/config"
	"orchestrator/pb"
	"sync"

	"google.golang.org/grpc"
)

var pingPongCliOnce sync.Once
var pingPongCli pb.PingPongServiceClient

func GetPingPongCliInstance() pb.PingPongServiceClient {
	burnnerCfg := config.GetConfigInstance().Client.Burnner
	pingPongCliOnce.Do(func() {

		targetGrpc := burnnerCfg.Host + ":" + burnnerCfg.Port
		conn, err := grpc.Dial(targetGrpc, grpc.WithInsecure())
		if err != nil {
			panic(err.Error())
		}
		pingPongCli = pb.NewPingPongServiceClient(conn)
	})
	return pingPongCli
}

package orchestrator

import (
	"fmt"
	"orchestrator/pkg/ctx"
	"testing"
)

// simulate main process
func init()  {
	var (
		handlerA, handlerB SyncNode
	)

	orchestrator := GetInstance()

	handlerA = func(requestID string, ctx *ctx.Context) error {
		// 將 request-id 寫入資料
		fmt.Println("handlerA", "request-id", requestID)
		ctx.Set("foo", "bar")
		return nil
	}

	handlerB = func(requestID string, ctx *ctx.Context) error{
		// 將 request-id 寫入資料
		fmt.Println("handlerB", "foo", ctx.Value("foo"))
		return nil
	}

	flow := NewSyncFlow()
	flow.Use(handlerA).Use(handlerB)
	orchestrator.SetSyncFlows(Facade("foo"), flow)
}

func Test_Sync(t *testing.T) {

	orchestrator := GetInstance()
	flow := orchestrator.GetSyncFlow(Facade("foo"))
	flow.Run("666", nil, "test-req", "test-resp")

}

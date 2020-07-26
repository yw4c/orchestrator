package orchestrator

import (
	"fmt"
	"orchestrator/pkg/ctx"
	"testing"
)

// simulate main process
func init()  {
	var (
		handlerA, handlerB SyncHandler
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

	flow := NewSyncFlow( Topic("666"))
	flow.Use(handlerA).Use(handlerB)
	orchestrator.SetFlows(Facade("foo"), flow)
}

func Test_Sync(t *testing.T) {

	orchestrator := GetInstance()
	flow := orchestrator.GetFlow(Facade("foo"))
	flow.Run("666")

}

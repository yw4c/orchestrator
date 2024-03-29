package orchestrator

import (
	"encoding/json"
	"github.com/rotisserie/eris"
	"orchestrator/pkg/pkgerror"
	"time"
)

const throttlingTimeout = 180

// 節流的事務流程
type IThrottlingFacade interface {
	// 將䩞點依序加入
	Use(TopicNodePair TopicNodePair) IAsyncFacade
	// 執行事務流程
	Run(requestID string, requestParam IAsyncFlowContext) (response interface{}, err error)
	IFacade
}

func NewThrottlingFlow(rollbackTopic Topic) *ThrottlingFlow {
	return &ThrottlingFlow{
		AsyncFacade: &AsyncFacade{
			pairs:         nil,
			rollbackTopic: rollbackTopic,
		},
	}
}

type ThrottlingFlow struct {
	*AsyncFacade
}

func (t *ThrottlingFlow) Run(requestID string, requestParam IAsyncFlowContext) (response interface{}, err error) {
	if len(t.pairs) == 0 {
		return
	}

	// 蒐集 topics
	var topics []Topic
	for _, v := range t.pairs {
		topics = append(topics, v.Topic)
	}

	// 準備 msg 給第一個事務
	requestParam.SetCurrentIndex(0)
	requestParam.SetRequestID(requestID)
	requestParam.SetTopics(topics)
	requestParam.SetRollbackTopic(t.rollbackTopic)

	data, err := json.Marshal(requestParam)
	if err != nil {
		return nil, eris.Wrap(pkgerror.ErrInternalError, "Json Marshal Fail")
	}

	// 開始推播給第一個事務
	GetMQInstance().Produce(t.pairs[0].Topic, data)
	response, err = Wait(requestID, throttlingTimeout*time.Second)
	return

}

package orchestrator

import (
	"encoding/json"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/pkg/pkgerror"
)

type IAsyncFlowContext interface {
	SetRequestID(requestID string)
	SetCurrentIndex(currentIndex int)
	SetTopics(topics []Topic)
	SetRollbackTopic(topic Topic)

	// 發布 topic 給下一個 Handler
	Next()
	// 打印 log ，發送 topic 給 rollback Handler
	Rollback(err error)
}

// Rollback 的推播格式
type RollbackMsg struct {
	RequestID string `json:"request_id"`
}
// 異步流程的推播格式
type AsyncFlowContext struct {
	// rollback topic
	RollbackTopic Topic `json:"rollback_topic"`
	// 請求編號
	RequestID string `json:"request_id"`
	// 當前節點索引 Topics，用於找到下個節點推播
	CurrentIndex int `json:"current_index"`
	// 流程中所有 topic (依序)
	Topics []Topic `json:"topics"`
}

func (a *AsyncFlowContext) Next() {
	var nextTopic Topic
	if len(a.Topics)-1 > a.CurrentIndex {
		nextTopic = a.Topics[a.CurrentIndex+1]
		a.CurrentIndex += 1
		nextData, _ := json.Marshal(a)
		GetMQInstance().Produce(nextTopic, nextData)
	}
}

func (a *AsyncFlowContext) Rollback(err error) {

	formattedStr := eris.ToCustomString(err, pkgerror.Format)

	log.Error().
		Str("track", formattedStr).
		Str("x-request-id", a.RequestID).
		Str("error", err.Error()).
		Msg("Consumer Error, Going to Rollback")

	rollback(a.RollbackTopic, a.RequestID)
}

func (a *AsyncFlowContext) SetRollbackTopic(topic Topic) {
	a.RollbackTopic = topic
}

func (a *AsyncFlowContext) SetRequestID( requestID string) {
	a.RequestID = requestID
}

func (a *AsyncFlowContext) SetCurrentIndex(currentIndex int) {
	a.CurrentIndex = currentIndex
}

func (a *AsyncFlowContext) SetTopics(topics []Topic) {
	a.Topics = topics
}

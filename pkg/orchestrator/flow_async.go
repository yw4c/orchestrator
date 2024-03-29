package orchestrator

import (
	"encoding/json"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/pkg/pkgerror"
	"runtime/debug"
)

// 異步的事務流程
type IAsyncFacade interface {
	// 註冊 Topic 與對應的 事務節點(AsyncNode)，有順序性
	Use(TopicNodePair TopicNodePair) IAsyncFacade
	// 開始準備接收 MQ 訊息
	Consume()
	// 執行事務流程
	Run(requestID string, requestParam IAsyncFlowContext) (err error)

	IFacade
}

// 異步的事務節點
type AsyncNode func(topic Topic, data []byte, next Next, rollback Rollback)

// 進行下一個事務節點
type Next func(context IAsyncFlowContext)

// 執行 Rollback：記錄錯誤日誌、推送 Rollback Topic
type Rollback func(err error, context IAsyncFlowContext)

// Rollback的事務節點
type RollbackNode func(topic Topic, data []byte)

// 一個 Topic 對應 一個 事務節點
type TopicNodePair struct {
	Topic     Topic
	AsyncNode AsyncNode
}
type TopicRollbackNodePair struct {
	Topic Topic
	Node  RollbackNode
}

// Push rollback topic
func rollback(topic Topic, requestID string) {
	msg, _ := json.Marshal(&RollbackMsg{RequestID: requestID})
	GetMQInstance().Produce(topic, msg)
}

type AsyncFacade struct {
	// 一個 Topic 對應一個事務
	pairs         []TopicNodePair
	rollbackTopic Topic
}

func (s *AsyncFacade) ConsumeRollback(rollback *TopicRollbackNodePair) {
	GetMQInstance().ConsumeRollback(rollback.Topic, rollback.Node)
	s.rollbackTopic = rollback.Topic
}

func (s *AsyncFacade) Consume() {
	if len(s.pairs) == 0 {
		return
	}
	mq := GetMQInstance()
	for _, v := range s.pairs {
		mq.ListenAndConsume(v.Topic, v.AsyncNode)
	}
}

func (s *AsyncFacade) Use(TopicNodePair TopicNodePair) IAsyncFacade {
	s.pairs = append(s.pairs, TopicNodePair)
	return s
}

func (s *AsyncFacade) Run(requestID string, requestParam IAsyncFlowContext) (err error) {
	if len(s.pairs) == 0 {
		return
	}

	// 蒐集 topics
	var topics []Topic
	for _, v := range s.pairs {
		topics = append(topics, v.Topic)
	}

	// 準備 msg 給第一個事務
	requestParam.SetCurrentIndex(0)
	requestParam.SetRequestID(requestID)
	requestParam.SetTopics(topics)
	requestParam.SetRollbackTopic(s.rollbackTopic)

	data, err := json.Marshal(requestParam)
	if err != nil {
		return eris.Wrap(pkgerror.ErrInternalError, "Json Marshal Fail")
	}

	// 開始推播給第一個事務
	GetMQInstance().Produce(s.pairs[0].Topic, data)
	return nil
}

func NewAsyncFacade(rollbackTopic Topic) *AsyncFacade {
	return &AsyncFacade{
		pairs:         nil,
		rollbackTopic: rollbackTopic,
	}
}

func getNextFunc() Next {
	return func(d IAsyncFlowContext) {

		var nextTopic Topic

		// has next？有的話 produce 下一個 topic
		if len(d.GetTopics())-1 > d.GetCurrentIndex() {

			nextTopic = d.GetTopics()[d.GetCurrentIndex()+1]
			d.SetCurrentIndex(d.GetCurrentIndex() + 1)
			nextData, err := json.Marshal(d)
			if err != nil {
				log.Error().Str("track", string(debug.Stack())).
					Interface("context", d).
					Msg("GetNextFunc Error")
				rollback(d.GetRollbackTopic(), d.GetRequestID())
				return
			}
			GetMQInstance().Produce(nextTopic, nextData)
		} else
		// It's the last one
		{
			// Throttling 模式, 通知 reqwait 任務已完成。並傳回 Context DTO
			currentTopic := d.GetTopics()[d.GetCurrentIndex()]
			if currentTopic.GetIsThrottling() {
				TaskFinished(d.GetRequestID(), d)
			}
		}
	}
}

func getRollbackFunc() Rollback {
	return func(err error, context IAsyncFlowContext) {
		// log tracked errors
		formattedStr := eris.ToCustomString(err, pkgerror.Format)
		log.Error().
			Str("track", formattedStr).
			Str("x-request-id", context.GetRequestID()).
			Str("error", err.Error()).
			Msgf("Trigger Rollback : %s", context.GetRollbackTopic())

		// Throttling 模式, 通知 reqwait 任務已完成。並傳回 error
		currentTopic := context.GetTopics()[context.GetCurrentIndex()]
		if currentTopic.GetIsThrottling() {
			TaskFinished(context.GetRequestID(), err)
		}

		// trigger rollback
		rollback(context.GetRollbackTopic(), context.GetRequestID())
	}
}

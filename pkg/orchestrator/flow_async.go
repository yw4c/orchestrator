package orchestrator

import (
	"encoding/json"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/pkg/pkgerror"
	"runtime/debug"
)

// 異步的事務流程
type IAsyncFlow interface {
	// 註冊 Topic 與對應的 事務節點(AsyncHandler)，有順序性
	Use(TopicHandlerPair TopicHandlerPair) IAsyncFlow
	// 開始準備接收 MQ 訊息
	Consume()
	// 執行事務流程
	Run(requestID string, requestParam IAsyncFlowMsg) (err error)

	IFlow
}

// 異步的事務節點
type AsyncHandler func(topic Topic, data []byte, next Next)
// 異步事務，進行下一個事務節點
type Next func(context IAsyncFlowMsg)
// Rollback的事務節點
type RollbackHandler func(topic Topic, data []byte)

// 一個 Topic 對應 一個 事務節點
type TopicHandlerPair struct {
	Topic        Topic
	AsyncHandler AsyncHandler
}
type TopicRollbackHandlerPair struct {
	Topic        Topic
	Handler RollbackHandler
}


// 推送 rollback Topic 給 mq
func rollback(topic Topic, requestID string) {
	msg, _ := json.Marshal(&RollbackMsg{RequestID: requestID})
	GetMQInstance().Produce(topic, msg)
}

type AsyncFlow struct {
	// 一個 Topic 對應一個事務
	handlers      []TopicHandlerPair
	rollbackTopic Topic
}

func (s *AsyncFlow) ConsumeRollback(rollback *TopicRollbackHandlerPair) {
	GetMQInstance().ConsumeRollback(rollback.Topic, rollback.Handler)
	s.rollbackTopic = rollback.Topic
}

func (s *AsyncFlow) Consume() {
	if len(s.handlers) == 0 {
		return
	}
	mq := GetMQInstance()
	for _, v := range s.handlers {
		mq.ListenAndConsume(v.Topic, v.AsyncHandler)
	}
}

func (s *AsyncFlow) Use(TopicHandlerPair TopicHandlerPair) IAsyncFlow {
	s.handlers = append(s.handlers, TopicHandlerPair)
	return s
}

func (s *AsyncFlow) Run(requestID string, requestParam IAsyncFlowMsg) (err error){
	if len(s.handlers) == 0 {
		return
	}

	// 蒐集 topics
	var topics []Topic
	for _, v := range s.handlers {
		topics = append(topics, v.Topic)
	}

	// 準備 msg 給第一個事務
	requestParam.SetCurrentIndex(0)
	requestParam.SetRequestID(requestID)
	requestParam.SetTopics(topics)
	requestParam.SetRollbackTopic(s.rollbackTopic)


	data, err := json.Marshal(requestParam)
	if err != nil {
		return  eris.Wrap(pkgerror.ErrInternalError, "Json Marshal Fail")
	}

	// 開始推播給第一個事務
	GetMQInstance().Produce(s.handlers[0].Topic, data)
	return    nil
}

func NewAsyncFlow(rollbackTopic Topic) *AsyncFlow {
	return &AsyncFlow{
		handlers: nil,
		rollbackTopic:rollbackTopic,
	}
}


func GetNextFunc() Next{
	return func(d IAsyncFlowMsg) {

		var nextTopic Topic

		// has next？有的話 produce 下一個 topic
		if len(d.GetTopics())-1 > d.GetCurrentIndex() {

			nextTopic = d.GetTopics()[d.GetCurrentIndex()+1]
			d.SetCurrentIndex(d.GetCurrentIndex()+1)
			nextData, err := json.Marshal(d)
			if err != nil {
				log.Error().Str("track", string(debug.Stack())).
					Interface("context", d).
					Msg("GetNextFunc Error")
				return
			}
			GetMQInstance().Produce(nextTopic, nextData)
		}
	}
}
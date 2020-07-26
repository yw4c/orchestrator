package orchestrator

import (
	"encoding/json"
)

// 異步的事務流程
type IAsyncFlow interface {
	// 註冊 Topic 與對應的 事務節點(AsyncHandler)，有順序性
	Use(TopicHandlerPair TopicHandlerPair) IAsyncFlow
	// 開始準備接收 MQ 訊息
	Consume()
	IFlow
}

// 異步的事務節點
type AsyncHandler func(topic Topic, data []byte)
// 取消所有 Flow 事務
//type RollbackHandler func(topic Topic, data []byte)
// 一個 Topic 對應 一個 事務節點
type TopicHandlerPair struct {
	Topic        Topic
	AsyncHandler AsyncHandler
}


// Rollback 的推播格式
type baseMsg struct {
	requestID string `json:"request_id"`
}
// 異步流程的推播格式
type AsyncFlowMsg struct {
	// 流程名稱
	Facade string `json:"facade"`
	// 請求編號
	RequestID string `json:"request_id"`
	// 當前節點索引 Topics，用於找到下個節點推播
	CurrentIndex int `json:"current_index"`
	// 流程中所有 topic (依序)
	Topics []string `json:"topics"`
}

// 推送 rollback Topic 給 mq
func rollback(topic Topic, requestID string) {
	msg, _ := json.Marshal(&baseMsg{requestID: requestID})
	GetMQInstance().Produce(topic, msg)
}

type AsyncFlow struct {
	// 一個 Topic 對應一個事務
	handlers      []TopicHandlerPair
	rollbackTopic Topic
}

func (s *AsyncFlow) ConsumeRollback(rollback *TopicHandlerPair) {
	panic("implement me")
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

func (s *AsyncFlow) Run(requestID string, requestParam interface{}) (response interface{}, err error) {
	if len(s.handlers) == 0 {
		return
	}

	// 蒐集 topics
	var topics []string
	for _, v := range s.handlers {
		topics = append(topics, string(v.Topic))
	}

	// 準備 msg 給第一個事務
	data, _ := json.Marshal(&AsyncFlowMsg{
		RequestID:requestID,
		CurrentIndex:0,
		Topics: topics,
	})

	// 開始推播給第一個事務
	GetMQInstance().Produce(s.handlers[0].Topic, data)
	return  requestID, nil
}

func NewAsyncFlow(rollbackTopic Topic) *AsyncFlow {
	return &AsyncFlow{
		handlers: nil,
		rollbackTopic:rollbackTopic,
	}
}

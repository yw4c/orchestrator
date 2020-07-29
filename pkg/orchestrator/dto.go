package orchestrator

type IAsyncFlowContext interface {
	SetRequestID(requestID string)
	SetCurrentIndex(currentIndex int)
	SetTopics(topics []Topic)
	SetRollbackTopic(topic Topic)

	GetRequestID()string
	GetCurrentIndex()int
	GetTopics()[]Topic
	GetRollbackTopic()Topic
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

func (a *AsyncFlowContext) GetRequestID() string {
	return a.RequestID
}

func (a *AsyncFlowContext) GetCurrentIndex() int {
	return a.CurrentIndex
}

func (a *AsyncFlowContext) GetTopics() []Topic {
	return a.Topics
}

func (a *AsyncFlowContext) GetRollbackTopic() Topic {
	return a.RollbackTopic
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

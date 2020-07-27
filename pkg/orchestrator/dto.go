package orchestrator

type IAsyncFlowMsg interface {
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
type AsyncFlowMsg struct {
	// rollback topic
	RollbackTopic Topic `json:"rollback_topic"`
	// 請求編號
	RequestID string `json:"request_id"`
	// 當前節點索引 Topics，用於找到下個節點推播
	CurrentIndex int `json:"current_index"`
	// 流程中所有 topic (依序)
	Topics []Topic `json:"topics"`
}

func (a *AsyncFlowMsg) GetRequestID() string {
	return a.RequestID
}

func (a *AsyncFlowMsg) GetCurrentIndex() int {
	return a.CurrentIndex
}

func (a *AsyncFlowMsg) GetTopics() []Topic {
	return a.Topics
}

func (a *AsyncFlowMsg) GetRollbackTopic() Topic {
	return a.RollbackTopic
}

func (a *AsyncFlowMsg) SetRollbackTopic(topic Topic) {
	a.RollbackTopic = topic
}

func (a *AsyncFlowMsg) SetRequestID( requestID string) {
	a.RequestID = requestID
}

func (a *AsyncFlowMsg) SetCurrentIndex(currentIndex int) {
	a.CurrentIndex = currentIndex
}

func (a *AsyncFlowMsg) SetTopics(topics []Topic) {
	a.Topics = topics
}

package orchestrator

type IAsyncFlowContext interface {
	SetRequestID(requestID string)
	SetCurrentIndex(currentIndex int)
	SetTopics(topics []Topic)
	SetRollbackTopic(topic Topic)

	GetRequestID() string
	GetCurrentIndex() int
	GetTopics() []Topic
	GetRollbackTopic() Topic
}

// Rollback format
type RollbackMsg struct {
	RequestID string `json:"request_id"`
}

// AsyncFacadeContext
type AsyncFacadeContext struct {
	// rollback topic
	RollbackTopic Topic `json:"rollback_topic"`
	// RequestID
	RequestID string `json:"request_id"`
	// Get current index of iterated nodes
	CurrentIndex int `json:"current_index"`
	// topics in order
	Topics []Topic `json:"topics"`
}

func (a *AsyncFacadeContext) GetRequestID() string {
	return a.RequestID
}

func (a *AsyncFacadeContext) GetCurrentIndex() int {
	return a.CurrentIndex
}

func (a *AsyncFacadeContext) GetTopics() []Topic {
	return a.Topics
}

func (a *AsyncFacadeContext) GetRollbackTopic() Topic {
	return a.RollbackTopic
}

func (a *AsyncFacadeContext) SetRollbackTopic(topic Topic) {
	a.RollbackTopic = topic
}

func (a *AsyncFacadeContext) SetRequestID(requestID string) {
	a.RequestID = requestID
}

func (a *AsyncFacadeContext) SetCurrentIndex(currentIndex int) {
	a.CurrentIndex = currentIndex
}

func (a *AsyncFacadeContext) SetTopics(topics []Topic) {
	a.Topics = topics
}

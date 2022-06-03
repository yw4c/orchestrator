package orchestrator

import (
	"orchestrator/config"
	"os"
	"sync"
)

type MQ interface {
	Produce(topic Topic, message []byte) error
	ListenAndConsume(topic Topic, node AsyncNode)
	ConsumeRollback(topic Topic, node RollbackNode)
}

var mq MQ
var onceMQ sync.Once

func GetMQInstance() MQ {
	onceMQ.Do(func() {
		mq = NewNats()
	})
	return mq
}

type Topic string

func (id Topic) GetTopicName() string {
	prefix := config.GetConfigInstance().MessageQueue.TopicPrefix + "."

	var name string
	for _, v := range config.GetConfigInstance().MessageQueue.Topics {
		if string(id) == v.ID {
			name = v.ID
			if v.IsThrottling {
				name = os.ExpandEnv(name + ".$POD_ID")
			}
			return prefix + name
		}
	}
	return prefix + string(id)
}

func (id Topic) GetConcurrency() int {
	for _, v := range config.GetConfigInstance().MessageQueue.Topics {
		if string(id) == v.ID {
			return v.Concurrency
		}
	}
	return 1
}

func (id Topic) GetIsThrottling() bool {
	for _, v := range config.GetConfigInstance().MessageQueue.Topics {
		if string(id) == v.ID {
			return v.IsThrottling
		}
	}
	return false
}

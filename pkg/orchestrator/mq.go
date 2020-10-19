package orchestrator

import (
	"orchestrator/config"
	"os"
	"sync"
)


type MQ interface {
	Produce(topic Topic, message []byte)
	ListenAndConsume(topic Topic, handler AsyncHandler)
	ConsumeRollback(topic Topic, handler RollbackHandler)
}

var mq MQ
var onceMQ sync.Once

func GetMQInstance() MQ {
	onceMQ.Do(func() {
		mq = NewRabbitMQ()
	})
	return mq
}

type Topic string

func (id Topic) GetTopicName() string {
	for _,v := range config.GetConfigInstance().Topics {
		if string(id) == v.ID {
			return os.ExpandEnv(v.Topic)
		}
	}
	return string(id)
}

func (id Topic) GetConcurrency() int {
	for _,v := range config.GetConfigInstance().Topics {
		if string(id) == v.ID {
			return v.Concurrency
		}
	}
	return 1
}
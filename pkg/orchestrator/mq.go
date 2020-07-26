package orchestrator

import (
	"sync"
)

type Topic string


type MQ interface {
	Produce(topic Topic, message []byte)
	ListenAndConsume(topic Topic, handler AsyncHandler)
}

var mq MQ
var onceMQ sync.Once

func GetMQInstance() MQ {
	onceMQ.Do(func() {
		mq = NewRabbitMQ()
	})
	return mq
}
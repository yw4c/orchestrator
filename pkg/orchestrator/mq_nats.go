package orchestrator

import (
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	"orchestrator/config"
	"runtime/debug"
	"strconv"
	"time"
)

var (
	NatsLostConn stan.Option
)

func NewNats() *Nats  {
	// 和 server 斷線的話
	NatsLostConn = stan.SetConnectionLostHandler(func(c stan.Conn, err error) {
		log.Error().Err(err).Msg("server dead")
		// 重新連線
		c.Close()
		connectNats(NatsLostConn)
		// 重新註冊
		//registerSubscribers()
	})

	return &Nats{
		conn: connectNats(NatsLostConn),
	}
}

type Nats struct {
	conn stan.Conn
}

func (n *Nats) Produce(topicID Topic, message []byte) {
	topic := topicID.GetTopicName()

	log.Info().Str("topic", topic).Str("msg", string(message)).Msg("Producing Msg")
	err := n.conn.Publish(topic, message)

	if err != nil {
		panic(err.Error())
	}
}

func (n *Nats) ListenAndConsume(topicID Topic, node AsyncNode) {
	topic := topicID.GetTopicName()


	for i:=1; i <= topicID.GetConcurrency();i ++ {
		consumerName := topic+"_consumer_"+strconv.Itoa(i)
		log.Info().Str("consumer_name", consumerName).Msg("registering")


		go func(node AsyncNode,consumerName string) {

			defer func() {
				if p := recover(); p != nil {
					log.Error().Interface("message", p).Str("trace", string(debug.Stack())).Msg("Consumer Panic Recover")
				}
			}()

			_, err := n.conn.QueueSubscribe(topic, "ackTimeout",
				func(msg *stan.Msg) {
					log.Info().
						Str("consumer", consumerName).
						Str("topic", topic).
						Str("body", string(msg.Data)).
						Msg("Consumer Access Log")

					node(topicID, msg.Data, getNextFunc(), getRollbackFunc())

					msg.Ack()
				},
				stan.SetManualAckMode(),
				stan.DurableName("ackTimeout"),
			)

			if err != nil {
				panic(err.Error())
			}
		}(node, consumerName)


	}
}

func (n *Nats) ConsumeRollback(topicID Topic, node RollbackNode) {
	topic := topicID.GetTopicName()


	for i:=1; i <= topicID.GetConcurrency();i ++ {
		consumerName := topic+"_consumer_"+strconv.Itoa(i)
		log.Info().Str("consumer_name", consumerName).Msg("registering")

		go func(node RollbackNode,consumerName string) {

			defer func() {
				if p := recover(); p != nil {
					log.Error().Interface("message", p).Str("trace", string(debug.Stack())).Msg("Consumer Panic Recover")
				}
			}()

			_, err := n.conn.QueueSubscribe(topic, "ackTimeout",
				func(msg *stan.Msg) {
					log.Info().
						Str("consumer", consumerName).
						Str("topic", topic).
						Str("body", string(msg.Data)).
						Msg("Consumer Access Log")

					node(topicID, msg.Data)

					msg.Ack()
				},
				stan.SetManualAckMode(),
				stan.DurableName("ackTimeout"),
			)

			if err != nil {
				panic(err.Error())
			}
		}(node, consumerName)


	}
}

func connectNats(opt ...stan.Option) (conn stan.Conn){

	natsCnf := config.GetConfigInstance().Client.Nats
	natsUrl := natsCnf.NatsUrl
	clusterID := natsCnf.ClusterId
	clientID := natsCnf.ClientId

	opt = append(opt, stan.NatsURL(natsUrl))
	var err error

	for {
		conn, err = stan.Connect(clusterID, clientID, opt...)

		// Reconnection
		if err != nil {
			log.Error().Err(err).Msg("connectNats failed")
		} else {
			log.Info().Msg("connectNats success")
			break
		}
		time.Sleep(time.Second)
	}
	return conn

}
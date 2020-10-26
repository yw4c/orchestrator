package orchestrator

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"orchestrator/config"
	"runtime/debug"
	"strconv"
	"time"
)

const(
	durable = true
	// 首次声明它的连接（Connection）可见
	rabbitmqExclusive = false
)

// NewRabbitMQ 建立 pod 與 RabbitMQ 連線
func NewRabbitMQ()*RabbitMQ {

	mqConfig := config.GetConfigInstance().Client.RabbitMQ
	mq := &RabbitMQ{}

	connQuery:= fmt.Sprintf("amqp://%s:%s@%s/", mqConfig.Username, mqConfig.Password, mqConfig.Host+":"+mqConfig.Port)
	log.Info().Msg("Dialing to RabbitMQ" + connQuery)

	var err error

	for {
		c := amqp.Config{
			ChannelMax: 10000,
		}
		mq.conn, err = amqp.DialConfig(connQuery, c)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
		log.Warn().Msg("rabbit mq dial retrying" )
	}

	return mq
}
type RabbitMQ struct {
	conn *amqp.Connection
}



func (r *RabbitMQ) Produce(topicID Topic, message []byte) {

	topic := topicID.GetTopicName()

	ch, err := r.conn.Channel()
	if err != nil {
		panic(err.Error())
	}
	//defer ch.Close()

	log.Info().Str("topic", string(topic)).Msg("Producing Msg")

	err = ch.Publish(
		string(topic+"_exchange"),     // exchange
		string(topic), // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	if err != nil {
		panic(err.Error())
	}
}


// Exchange -> Binding -> Queues -> Consumer
func (r *RabbitMQ) ListenAndConsume(topicID Topic, node AsyncNode) {

	topic := topicID.GetTopicName()

	// channels 復用一 Connection 連線, channel 對應一 consumer
	channel, err := r.conn.Channel()
	if err != nil {
		panic(err.Error())
	}



	// 使用 direct 一對一接收
	err = channel.ExchangeDeclare(
		string(topic+"_exchange"),   // name
		"direct", // type
		false,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		panic(err.Error())
	}


	// Queue
	q, err := channel.QueueDeclare(
		string(topic+"_queue"),    // name
		durable, // durable
		false, // delete when unused
		rabbitmqExclusive,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err.Error())
	}


	// Bindings
	err = channel.QueueBind(
		q.Name, // queue name, 这里指的是 test_logs
		string(topic),     // routing key
		string(topic+"_exchange"), // exchange
		false,
		nil)
	if err != nil {
		panic(err.Error())
	}

	for i:=1; i <= topicID.GetConcurrency();i ++ {
		consumerName := topic+"_consumer_"+strconv.Itoa(i)
		log.Info().Str("consumer_name", consumerName).Msg("registering")

		// register consumer
		msgs, err := channel.Consume(
			q.Name,                    // queue
			consumerName, // consumer
			false,                     // auto-ack
			rabbitmqExclusive,                     // exclusive
			false,                     // no-local
			false,                     // no-wait
			nil,                       // args
		)
		if err != nil {
			panic(err.Error())
		}

		go func(msgs <-chan amqp.Delivery, node AsyncNode, consumerName string) {

			// Panic recover
			defer func() {
				if p := recover(); p != nil {
					log.Error().Interface("message", p).Str("trace", string(debug.Stack())).Msg("Consumer Panic Recover")
				}
				defer r.conn.Close()
				defer channel.Close()
			}()

			// 處理訊息
			for d := range msgs {
				log.Info().
					Str("consumer", consumerName).
					Str("topic", string(topic)).
					Str("body", string(d.Body)).
					Msg("Consumer Access Log")

				node(topicID, d.Body, getNextFunc(), getRollbackFunc())

				d.Ack(false)

			}
		}(msgs, node, consumerName)

	}

}

func (r *RabbitMQ) ConsumeRollback(topicID Topic, node RollbackNode) {

	topic := topicID.GetTopicName()

	// channels 復用一 Connection 連線, channel 對應一 consumer
	channel, err := r.conn.Channel()
	if err != nil {
		panic(err.Error())
	}



	// 使用 direct 一對一接收
	err = channel.ExchangeDeclare(
		string(topic+"_exchange"),   // name
		"direct", // type
		false,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		panic(err.Error())
	}


	// Queue
	q, err := channel.QueueDeclare(
		string(topic+"_queue"),    // name
		true, // durable
		false, // delete when unused
		rabbitmqExclusive,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err.Error())
	}


	// Bindings
	err = channel.QueueBind(
		q.Name, // queue name, 这里指的是 test_logs
		string(topic),     // routing key
		string(topic+"_exchange"), // exchange
		false,
		nil)
	if err != nil {
		panic(err.Error())
	}

	for i:=1; i <= topicID.GetConcurrency();i ++ {
		consumerName := topic + "_consumer_" + strconv.Itoa(i)
		log.Info().Str("consumer_name", consumerName).Msg("registering")

		// register consumer
		msgs, err := channel.Consume(
			q.Name,                    // queue
			consumerName, // consumer
			false,                     // auto-ack
			rabbitmqExclusive,                     // exclusive
			false,                     // no-local
			false,                     // no-wait
			nil,                       // args
		)
		if err != nil {
			panic(err.Error())
		}

		go func(msgs <-chan amqp.Delivery, node RollbackNode) {

			// Panic recover
			defer func() {
				if p := recover(); p != nil {
					log.Error().Interface("message", p).Str("trace", string(debug.Stack())).Msg("Consumer Panic Recover")
				}
				defer r.conn.Close()
				defer channel.Close()
			}()

			// 處理訊息
			for d := range msgs {
				log.Info().
					Str("topic", string(topic)).
					Str("body", string(d.Body)).
					Msg("Consumer Access Log")

				node(topicID, d.Body)

				d.Ack(false)
			}
		}(msgs, node)
	}

}
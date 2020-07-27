package orchestrator

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"runtime/debug"
	"time"
)


// NewRabbitMQ 建立 pod 與 RabbitMQ 連線
func NewRabbitMQ()*RabbitMQ {

	mq := &RabbitMQ{}

	connQuery:= fmt.Sprintf("amqp://%s:%s@%s/","guest" , "guest", "localhost:5672")
	log.Info().Msg("Dialing to RabbitMQ" + connQuery)

	var err error

	for {
		mq.conn, err = amqp.Dial(connQuery)
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


func (r *RabbitMQ) Produce(topic Topic, message []byte) {
	ch, err := r.conn.Channel()
	if err != nil {
		panic(err.Error())
	}
	defer ch.Close()

	if topic == "" {
		panic("topic can not be empty")
	}

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
func (r *RabbitMQ) ListenAndConsume(topic Topic, handler AsyncHandler) {



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
		true,  // exclusive
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


	// register consumer
	msgs, err := channel.Consume(
		q.Name, // queue
		string(topic+"_consumer"),     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err.Error())
	}

	go func(msgs <-chan amqp.Delivery, handler AsyncHandler) {

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

			handler(topic, d.Body)

			d.Ack(false)
		}
	}(msgs, handler)



}


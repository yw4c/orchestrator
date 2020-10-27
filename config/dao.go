package config

type config struct {
	//  產品環境 <dev,sit,prod>
	Env string
	IsDebug bool `mapstructure:"is_debug"`
	// app 對外服務
	Server struct{
		// gRPC
		Grpc struct{
			Port string
		}
	}
	// 抽換服務
	Engine struct{
		// message queue 選用 <rabbit_mq, nats>
		MessageQueue string  `mapstructure:"message_queue"`
	}
	// Message Queue Topics
	MessageQueue struct{
		TopicPrefix string `mapstructure:"topic_prefix"`
		Topics []Topic
	} `mapstructure:"message_queue"`

	// app 連線對象
	Client struct {
		Redis Redis
		RabbitMQ RabbitMQ `mapstructure:"rabbit_mq"`
		Nats Nats `mapstructure:"nats"`
	}

}

type Topic struct {
	// mq server 註冊的 topic 名稱
	IsThrottling bool `mapstructure:"is_throttling"`
	// unique name
	ID string `mapstructure:"id"`
	// 接收的併發數
	Concurrency int
}

type Nats struct {
	ClusterId string `mapstructure:"cluster_id"`
	ClientId string `mapstructure:"client_id"`
	NatsUrl string `mapstructure:"nats_url"`
}

type RabbitMQ struct {
	Host string
	Port string
	Username string
	Password string
}

type Redis struct {
	Host string
	Port string
	DB int
	PoolSize int
	MinIdleConn int
	DialTimeout int
	ReadTimeout int
	WriteTimeout int
	PoolTimeout int
	IdleTimeout int
}
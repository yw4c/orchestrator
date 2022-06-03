package config

type config struct {
	// Env dev,sit,prod
	Env     string
	IsDebug bool `mapstructure:"is_debug"`

	Server struct {
		Grpc struct {
			Port string
		}
	}
	Engine struct {
		// MessageQueue <rabbit_mq, nats>
		MessageQueue string `mapstructure:"message_queue"`
	}
	// Message Queue Topics
	MessageQueue struct {
		TopicPrefix string `mapstructure:"topic_prefix"`
		Topics      []Topic
	} `mapstructure:"message_queue"`

	Client struct {
		Redis    Redis
		RabbitMQ RabbitMQ `mapstructure:"rabbit_mq"`
		Nats     Nats     `mapstructure:"nats"`
		Burnner  Burnner  `mapstructure:"burnner-grpc"`
	}

	Throttling struct {
		Concurrency int
		Enable      bool
	}
}

type Topic struct {
	IsThrottling bool `mapstructure:"is_throttling"`
	// unique name
	ID string `mapstructure:"id"`
	// Count of concurrency
	Concurrency int
}

type Nats struct {
	ClusterId string `mapstructure:"cluster_id"`
	ClientId  string `mapstructure:"client_id"`
	NatsUrl   string `mapstructure:"nats_url"`
}

type RabbitMQ struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Redis struct {
	Host         string
	Port         string
	DB           int
	PoolSize     int
	MinIdleConn  int
	DialTimeout  int
	ReadTimeout  int
	WriteTimeout int
	PoolTimeout  int
	IdleTimeout  int
}

type Burnner struct {
	Host string
	Port string
}

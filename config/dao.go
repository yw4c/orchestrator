package config

type config struct {
	//  產品環境 <dev,sit,prod>
	Env string
	IsDebug bool
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
	Topics []Topic

	// app 連線對象
	Client struct {
		Redis Redis
		RabbitMQ RabbitMQ `mapstructure:"rabbit_mq"`
	}

}

type Topic struct {
	// short name
	Topic string
	// mq server 註冊的 topic 名稱
	Name string
	// 接收的併發數
	Concurrency int
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
}
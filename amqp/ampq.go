package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/conf"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName = "my_queue"
)

func InitAmqp() error {
	// 创建RabbitMQ连接
	var err error
	config := conf.AmqpConfig
	conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", config.User, config.Passwd, config.Host, config.Port))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return err
	}

	// 创建通道
	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return err
	}

	// 声明一个队列
	queueName := "my_queue"
	_, err = ch.QueueDeclare(
		queueName, // 队列名称
		false,     // 是否持久化
		false,     // 是否自动删除
		false,     // 是否独占模式
		false,     // 是否阻塞等待
		nil,       // 额外参数
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
		return err
	}
	return nil
}

func SendMessage(data interface{}) {
	// 序列化结构体为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal struct to JSON: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发布消息到队列
	err = ch.PublishWithContext(
		ctx,
		"",        // 交换器名称，使用默认交换器
		queueName, // 队列名称
		false,     // 是否立即发送消息
		false,     // 是否等待服务器确认
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	fmt.Println("Struct sent successfully.")
}

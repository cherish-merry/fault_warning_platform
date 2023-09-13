package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/conf"
	"github.com/RaymondCode/simple-demo/database"
	"github.com/RaymondCode/simple-demo/models"
	log "github.com/sirupsen/logrus"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn  *amqp.Connection
	ch    *amqp.Channel
	go2py = "go2py"
	py2go = "py2go"
)

func InitAmqp() error {
	// 创建RabbitMQ连接
	var err error
	config := conf.AmqpConfig
	conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", config.User, config.Passwd, config.Host, config.Port))
	if err != nil {
		log.Errorf("Failed to connect to RabbitMQ: %v", err)
		return err
	}

	// 创建通道
	ch, err = conn.Channel()
	if err != nil {
		log.Errorf("Failed to open a channel: %v", err)
		return err
	}

	_, err = ch.QueueDeclare(
		go2py, // 队列名称
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否独占模式
		false, // 是否阻塞等待
		nil,   // 额外参数
	)
	if err != nil {
		log.Errorf("Failed to declare a queue: %v", err)
		return err
	}
	return nil
}

func SendMessage(data interface{}) {
	// 序列化结构体为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Failed to marshal struct to JSON: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发布消息到队列
	err = ch.PublishWithContext(
		ctx,
		"",    // 交换器名称，使用默认交换器
		go2py, // 队列名称
		false, // 是否立即发送消息
		false, // 是否等待服务器确认
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		log.Errorf("Failed to publish a message: %v", err)
	}

	log.Info("Struct sent successfully.")
}

func HandlerMessage() {
	q, err := ch.QueueDeclare(
		py2go,
		false, // 持久化
		false, // 不自动删除
		false, // 不独占
		false, // 不阻塞
		nil,
	)
	// 获取队列中的消息
	feedbacks, err := ch.Consume(
		q.Name, // 队列名称
		"",     // 消费者标识符
		true,   // 自动应答
		false,  // 不独占
		false,  // 不阻塞
		false,  // 不等待服务器确认
		nil,
	)
	if err != nil {
		log.Infof("无法获取消息: %v", err)
	}

	// 处理收到的消息
	forever := make(chan bool)

	go func() {
		for msg := range feedbacks {
			feedback := models.Feedback{}
			err := json.Unmarshal(msg.Body, &feedback)
			if err != nil {
				log.Errorf("Unmarshal feedback fail: %v", err)
			}
			db := database.GetInstanceConnection().GetPrimaryDB()
			log.Infof("收到消息: %v\n", feedback)
			err = feedback.Create(db)
			if err != nil {
				log.Errorf("Create feedback fail: %v", err)
			}

			// 添加处理消息的逻辑
		}
	}()

	log.Info("等待消息。按 Ctrl+C 退出")
	<-forever
}

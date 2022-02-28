package rabbitmq

import (
	"bytes"
	"errors"
	"github.com/streadway/amqp"
)

// 创建simple简单模式实例
func (this *RabbitMQ) NewSimple(queueName string) error {
	return this.New(queueName, "", "")
}

// 申请队列,如果队列不存在会自动创建,如果存在则跳过创建,保证队列存在,消息队列能发送到队列中
func (this *RabbitMQ) applicationSimpleQueue() error {
	_, err := this.channel.QueueDeclare(
		this.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		return errors.New("applicationSimpleQueue Error: " + err.Error())
	}
	return nil
}

// 简单模式发送消息
func (this *RabbitMQ) PublishSimple(message string) error {
	err := this.applicationSimpleQueue()
	if err != nil {
		return err
	}
	err = this.channel.Publish(
		this.Exchange,
		this.QueueName,
		//如果为true,根据exchange类型和routekey规则,如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//如果为true,当exchange发送消息队列到队列后发现队列上没有绑定消费者,则会把消息发还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return errors.New("PublishSimple Error: " + err.Error())
	}
	return nil
}

// 简单模式消费消息
func (this *RabbitMQ) ConsumeSimple(reader func(msg *string)) error {
	err := this.applicationSimpleQueue()
	if err != nil {
		return err
	}

	// 接受消息
	msgs, err := this.channel.Consume(
		this.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true,表示不能将同一个connection中发送消息传递给这个connection中的消费者
		false,
		//队列消费是否阻塞
		false,
		nil)
	if err != nil {
		return errors.New("ConsumeSimple Error: " + err.Error())
	}

	// 启用协程处理消息
	go func() {
		for d := range msgs {
			r := bytes.NewBuffer(d.Body).String()
			reader(&r) //实现要处理的逻辑函数
		}
	}()

	// 需要接收消息时此方法应该持续运行，后台服务一般在出口处就将程序阻塞住，所以这里可以不需要处理
	// select {}
	return nil
}

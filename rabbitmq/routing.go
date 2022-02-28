package rabbitmq

import (
	"bytes"
	"errors"
	"github.com/streadway/amqp"
)

// 创建路由模式实例
func (this *RabbitMQ) NewRouting(exchangeName, routingKey string) error {
	return this.New("", exchangeName, routingKey)
}

// 尝试创建交换机
func (this *RabbitMQ) applicationRoutingExchange() error {
	err := this.channel.ExchangeDeclare(
		this.Exchange,
		//交换机类型
		"direct",
		true,
		false,
		//YES表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New("applicationRoutingExchange Error: " + err.Error())
	}
	return nil
}

// 路由模式发送消息
func (this *RabbitMQ) PublishRouting(message string) error {
	err := this.applicationRoutingExchange()
	if err != nil {
		return err
	}
	err = this.channel.Publish(
		this.Exchange,
		this.key,
		//如果为true,根据exchange类型和routekey规则,如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//如果为true,当exchange发送消息队列到队列后发现队列上没有绑定消费者,则会把消息发还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return errors.New("PublishPub Error: " + err.Error())
	}
	return nil
}

// 路由模式消费消息
func (this *RabbitMQ) ConsumeRouting(reader func(msg *string)) error {
	err := this.applicationRoutingExchange()
	if err != nil {
		return err
	}

	// 创建队列，这里注意队列名称不要写
	q, err := this.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 绑定队列到 exchange 中
	err = this.channel.QueueBind(
		q.Name,
		// 需要绑定key
		this.key,
		this.Exchange,
		false,
		nil,
	)

	// 接受消息
	msgs, err := this.channel.Consume(
		q.Name,
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
		return errors.New("ConsumeRouting Error: " + err.Error())
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

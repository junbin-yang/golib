package rabbitmq

import (
	"errors"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//连接信息
	Host  string
	Port  string
	Vhost string
	User  string
	Pass  string
	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	key string
}

//断开channel和connection
func (this *RabbitMQ) Close() {
	this.channel.Close()
	this.conn.Close()
}

//创建RabbitMQ结构体实例
func (this *RabbitMQ) New(queueName, exchange, key string) error {
	if this.Host == "" {
		this.Host = "mq.iptv.sunlight-tech.com"
	}
	if this.Port == "" {
		this.Port = "5566"
	}
	if this.User == "" {
		this.User = "sun"
	}
	if this.Pass == "" {
		this.Pass = "sunlight2010"
	}
	if this.Vhost == "" {
		this.Vhost = "dvs_noticesvr"
	}
	this.QueueName = queueName
	this.Exchange = exchange
	this.key = key
	Mqurl := "amqp://" + this.User + ":" + this.Pass + "@" + this.Host + ":" + this.Port + "/" + this.Vhost
	var err error
	this.conn, err = amqp.Dial(Mqurl)
	if err != nil {
		return errors.New("创建连接错误: " + err.Error())
	}
	this.channel, err = this.conn.Channel()
	if err != nil {
		return errors.New("获取channel失败: " + err.Error())
	}
	return nil
}

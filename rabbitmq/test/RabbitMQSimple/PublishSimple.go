package main

import (
	"fmt"
	MQ "github.com/junbin-yang/golib/rabbitmq"
)

// 简单模式发送MQ消息
func main() {
	rabbitmq := &MQ.RabbitMQ{Vhost: "noticesvr"}
	err := rabbitmq.NewSimple("testSimple")
	failOnErr(err)
	err = rabbitmq.PublishSimple("Hello test!")
	failOnErr(err)

	fmt.Println("发送成功!")
	rabbitmq.Close()
}

func failOnErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}

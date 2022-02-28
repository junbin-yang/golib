package main

import (
	"fmt"
	MQ "github.com/junbin-yang/golib/rabbitmq"
)

// 简单模式接收MQ消息处理
func main() {
	rabbitmq := &MQ.RabbitMQ{Vhost: "noticesvr"}
	err := rabbitmq.NewSimple("testSimple")
	failOnErr(err)
	err = rabbitmq.ConsumeSimple(func(msg *string) {
		fmt.Printf("receve msg is :%s\n", *msg)
		//实现要处理的其他逻辑...
	})
	failOnErr(err)
	select {}
}

func failOnErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}

package main

import (
	"fmt"
	MQ "golib/rabbitmq"
	"strconv"
	"time"
)

// 订阅模式发送MQ消息
func main() {
	rabbitmq := &MQ.RabbitMQ{Vhost: "noticesvr"}
	err := rabbitmq.NewPubSub("", "newProduct")
	failOnErr(err)

	for i := 0; i < 10; i++ {
		err = rabbitmq.PublishPub("订阅模式生产第" + strconv.Itoa(i) + "条数据")
		failOnErr(err)
		fmt.Println("发送订阅模式生产第" + strconv.Itoa(i) + "条数据")
		time.Sleep(1 * time.Second)
	}
	rabbitmq.Close()
}

func failOnErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}

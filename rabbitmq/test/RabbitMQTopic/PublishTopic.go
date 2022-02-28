package main

import (
	"fmt"
	MQ "github.com/junbin-yang/golib/rabbitmq"
	"strconv"
	"time"
)

// 话题模式发送MQ消息
func main() {
	rabbitmq := &MQ.RabbitMQ{Vhost: "noticesvr"}
	err := rabbitmq.NewTopic("topictest", "my.topic.one")
	failOnErr(err)

	rabbitmq2 := &MQ.RabbitMQ{Vhost: "noticesvr"}
	err = rabbitmq2.NewTopic("topictest", "my.topic.two")
	failOnErr(err)

	for i := 0; i < 10; i++ {
		rabbitmq.PublishTopic("Hello topic one!" + strconv.Itoa(i))
		rabbitmq2.PublishTopic("Hello topic Two!" + strconv.Itoa(i))
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
	rabbitmq.Close()
}

func failOnErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}

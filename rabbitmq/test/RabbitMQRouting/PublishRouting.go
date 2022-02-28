package main

import (
	"fmt"
	MQ "github.com/junbin-yang/golib/rabbitmq"
	"strconv"
	"time"
)

// 路由模式发送MQ消息
func main() {
	rabbitmq := &MQ.RabbitMQ{Vhost: "noticesvr"}
	err := rabbitmq.NewRouting("newRouting", "keyOne")
	failOnErr(err)

	rabbitmq2 := &MQ.RabbitMQ{Vhost: "noticesvr"}
	err = rabbitmq2.NewRouting("newRouting", "keyTwo")
	failOnErr(err)

	for i := 0; i < 10; i++ {
		rabbitmq.PublishRouting("keyOne: 路由模式生产第" + strconv.Itoa(i) + "条数据")
		rabbitmq2.PublishRouting("keyTwo: 路由模式生产第" + strconv.Itoa(i) + "条数据")
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

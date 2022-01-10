package main

import (
	"fmt"
	"golib/redisx"
	"math/rand"
	"time"
)

var cluster, alone redisx.Rediser

func main() {
	redisAlone()
	fmt.Println("")
	redisCluster()
}

func redisAlone() {
	alone = &redisx.Alone{Host:"127.0.0.1:6379",Pass:""}
	alone.Connect()
	for i := 0; i < 5; i++ {
		fmt.Println("======== Redis Alone Test ========", i)
		alone.Set("key1", RandString(32), 60)
		fmt.Println(alone.Get("key1"))
		time.Sleep(2 * time.Second)
	}
}

func redisCluster() {
	cluster = &redisx.Cluster{Host:[]string{"127.0.0.1:6379","127.0.0.1:6380"},Pass:""}
	cluster.Connect()
	for i := 0; i < 5; i++ {
		fmt.Println("======= Redis Cluster Test =======", i)
		cluster.Set("key1", RandString(32), 60)
		fmt.Println(cluster.Get("key1"))
		time.Sleep(2 * time.Second)
	}
}

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

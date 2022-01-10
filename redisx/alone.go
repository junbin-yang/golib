package redisx

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type Alone struct {
	Client *redis.Pool
	Host   string
	Pass   string
}

func (this *Alone) Connect() {
	if this.Host == "" {
		this.Host = "127.0.0.1:6379"
	}

	this.Client = &redis.Pool{
		MaxIdle:     100,             // 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭
		MaxActive:   1024,            // 最大连接数，即最多的tcp连接数
		IdleTimeout: time.Second * 5, // 空闲连接超时时间（超时会关闭，并释放可用连接数）
		Wait:        true,            // 如果超过最大连接，是报错，还是等待
		Dial: func() (redis.Conn, error) {
			if this.Pass == "" {
				return redis.Dial("tcp", this.Host)
			}
			return redis.Dial("tcp", this.Host, redis.DialPassword(this.Pass))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error { // 如果设置了给func,那么每次p.Get()的时候都会调用改方法来验证连接的可用性
			_, err := c.Do("PING")
			return err
		},
	}
}

func (this *Alone) Expire(key string, time int64) error {
	conn := this.Client.Get()
	defer conn.Close()

	if time > 0 {
		_, err := conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Alone) Set(key string, data interface{}, time int64) error {
	conn := this.Client.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, data)
	if err != nil {
		return err
	}

	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}

	return nil
}

// Exists check a key
func (this *Alone) Exists(key string) (bool, int64) {
	conn := this.Client.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, 0
	}

	var ttl int64
	if exists {
		ttl, _ = redis.Int64(conn.Do("TTL", key))
	}

	return exists, ttl
}

// Get get a key
func (this *Alone) Get(key string) (string, error) {
	conn := this.Client.Get()
	defer conn.Close()

	return redis.String(conn.Do("GET", key))
}

// Delete delete a kye
func (this *Alone) Delete(key string) (bool, error) {
	conn := this.Client.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func (this *Alone) LikeDeletes(key string) error {
	conn := this.Client.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = this.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Alone) Like(key string) ([]string, error) {
	conn := this.Client.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("KEYS", "*"+key+"*"))
}

func (this *Alone) Sadd(key string, data interface{}, time int64) error {
	conn := this.Client.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", key, data)
	if err != nil {
		return err
	}

	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}

	return nil
}

// 判断元素是否是集合的成员
func (this *Alone) Sismember(key string, data interface{}) bool {
	conn := this.Client.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("SISMEMBER", key, data))
	if err != nil {
		return false
	}

	return exists
}

func (this *Alone) SetNX(key string, data interface{}, time int64) bool {
	conn := this.Client.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("SETNX", key, data))
	if err != nil {
		return false
	}

	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return false
		}
	}

	return exists
}

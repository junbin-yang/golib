package redisx

import (
	"github.com/junbin-yang/redis-go-cluster"
	"log"
	"time"
)

type Cluster struct {
	Client *redis.Cluster
	Host   []string
	Pass   string
}

func (this *Cluster) Connect() {
	var err error
	this.Client, err = redis.NewCluster(
		&redis.Options{
			StartNodes:   this.Host,
			ConnTimeout:  50 * time.Millisecond,
			ReadTimeout:  50 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
			KeepAlive:    16,
			AliveTime:    60 * time.Second,
			Password:     this.Pass,
		},
	)
	if err != nil {
		log.Fatalf("redis.Cluster error: %s", err.Error())
	}
}

func (this *Cluster) Expire(key string, time int64) error {
	conn := this.Client

	if time > 0 {
		_, err := conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Cluster) Set(key string, data interface{}, time int64) error {
	conn := this.Client

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
func (this *Cluster) Exists(key string) (bool, int64) {
	conn := this.Client

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
func (this *Cluster) Get(key string) (string, error) {
	conn := this.Client

	return redis.String(conn.Do("GET", key))
}

// Delete delete a kye
func (this *Cluster) Delete(key string) (bool, error) {
	conn := this.Client

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func (this *Cluster) LikeDeletes(key string) error {
	conn := this.Client

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

func (this *Cluster) Like(key string) ([]string, error) {
	conn := this.Client
	return redis.Strings(conn.Do("KEYS", "*"+key+"*"))
}

func (this *Cluster) Sadd(key string, data interface{}, time int64) error {
	conn := this.Client

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
func (this *Cluster) Sismember(key string, data interface{}) bool {
	conn := this.Client

	exists, err := redis.Bool(conn.Do("SISMEMBER", key, data))
	if err != nil {
		return false
	}

	return exists
}

func (this *Cluster) SetNX(key string, data interface{}, time int64) bool {
	conn := this.Client

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

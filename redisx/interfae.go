package redisx

type Rediser interface {
	Connect()
	Exists(key string) (bool, int64)
	Expire(key string, time int64) error
	Get(key string) (string, error)
	Delete(key string) (bool, error)
	LikeDeletes(key string) error
	Like(key string) ([]string, error)
	Set(key string, data interface{}, time int64) error
	Sadd(key string, data interface{}, time int64) error
	Sismember(key string, data interface{}) bool
	SetNX(key string, data interface{}, time int64) bool
}

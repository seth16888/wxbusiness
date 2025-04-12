package redis

import (
	"context"

	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var once sync.Once

var Redis *RedisClient

type RedisClient struct {
	Client  *redis.Client
	Context context.Context
	logger  *zap.Logger
}

// ConnectRedis 连接Redis
func ConnectRedis(addr, username, password string, db int, log *zap.Logger) {
	once.Do(func() {
		Redis = NewClient(addr, username, password, db, log)
	})
}

func NewClient(addr, username, password string, db int, logger *zap.Logger) *RedisClient {
	rds := &RedisClient{
		logger: logger,
	}
	rds.Context = context.Background()

	rds.Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	// ping
	err := rds.Ping()
	if err != nil {
    logger.Error("Connect", zap.Error(err))
	}

	return rds
}

func (rds *RedisClient) Ping() error {
	return rds.Client.Ping(rds.Context).Err()
}

// Set 设置key和value, and expire time
func (rds *RedisClient) Set(key string, value interface{}, expiration time.Duration) bool {
	if err := rds.Client.Set(rds.Context, key, value, expiration).Err(); err != nil {
    rds.logger.Error("Redis Set error", zap.Error(err))
		return false
	}
	return true
}

// Get 获取key对应的value
func (rds *RedisClient) Get(key string) (interface{}, error) {
	if value, err := rds.Client.Get(rds.Context, key).Result(); err != nil {
		if err != redis.Nil {
      rds.logger.Error("get", zap.Error(err))
		}
		return nil, err
	} else {
		return value, nil
	}
}

// Has 判断key是否存在
func (rds *RedisClient) Has(key string) bool {
	if _, err := rds.Client.Exists(rds.Context, key).Result(); err != nil {
		if err != redis.Nil {
			rds.logger.Error("Has", zap.Error(err))
		}
		return false
	}
	return true
}

// Del 删除keys
func (rds *RedisClient) Del(keys ...string) error {
	if _, err := rds.Client.Del(rds.Context, keys...).Result(); err != nil {
		rds.logger.Error("Del", zap.Error(err))
		return err
	}
	return nil
}

// FlushDB 清空当前数据库
func (rds *RedisClient) FlushDB() error {
	if err := rds.Client.FlushDB(rds.Context).Err(); err != nil {
		rds.logger.Error("FlushDB", zap.Error(err))
		return err
	}
	return nil
}

// Incr 增加key对应的value
func (rds *RedisClient) Incr(key string) (int64, error) {
	if v, err := rds.Client.Incr(rds.Context, key).Result(); err != nil {
		rds.logger.Error("Incr", zap.Error(err))
		return 0, err
	} else {
		return v, nil
	}
}

// Decr 减少key对应的value
func (rds *RedisClient) Decr(key string) (int64, error) {
	if v, err := rds.Client.Decr(rds.Context, key).Result(); err != nil {
		rds.logger.Error("Decr", zap.Error(err))
		return 0, err
	} else {
		return v, nil
	}
}

// IncrBy 增加key对应的value
func (rds *RedisClient) IncrBy(key string, value int64) (int64, error) {
	if v, err := rds.Client.IncrBy(rds.Context, key, value).Result(); err != nil {
		rds.logger.Error("IncrBy", zap.Error(err))
		return 0, err
	} else {
		return v, nil
	}
}

// DecrBy 减少key对应的value
func (rds *RedisClient) DecrBy(key string, value int64) (int64, error) {
	if v, err := rds.Client.DecrBy(rds.Context, key, value).Result(); err != nil {
		rds.logger.Error("DecrBy", zap.Error(err))
		return 0, err
	} else {
		return v, nil
	}
}

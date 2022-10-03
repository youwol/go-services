package utils

// This file is just a cache wrapper for github.com/patrickmn/go-cache and github.com/go-redis/redis
// We use it as a convenience to be able to switch between local and external caches
import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

var redisHost = os.Getenv("REDIS_HOST_PORT")

// CacheType for wrapper initialization (0 - local, 1 - redis)
type CacheType int8

const (
	// Local is the enum for local in-memory store
	Local CacheType = iota + 1
	// Redis is the enum for Redis cache store
	Redis
)

// Cache defines the interface for the cache wrapper
type Cache interface {
	Get(ctx context.Context, key string, value interface{}) bool
	Set(ctx context.Context, key string, value interface{}) error
}

// LocalCache is a in-memory cache store
type LocalCache struct {
	Client   *cache.Cache
	Duration time.Duration
}

// RedisCache uses the Redis instance from the infra namespace
type RedisCache struct {
	Client   *redis.Client
	Duration time.Duration
	Prefix   string
}

var cacheWrapper = make(map[string]Cache)

// InitCache intializes a cache object of the given type
func InitCache(name string, ct CacheType, d time.Duration) error {
	switch ct {
	case Local:
		cacheWrapper[name] = LocalCache{Client: cache.New(d, d+1*time.Second), Duration: d}
	case Redis:
		cacheWrapper[name] = RedisCache{
			Client: redis.NewClient(&redis.Options{
				Addr:     redisHost,
				Password: "", // no password set
				DB:       0,  // use default DB
			}),
			Duration: d,
			Prefix:   name,
		}
	default:
		return fmt.Errorf("Unkown cache type")
	}

	return nil
}

// GetCache retrieves a value from the cache store
func GetCache(ctx context.Context, name string, key string, value interface{}) bool {
	c, found := cacheWrapper[name]
	logger := ContextLogger(ctx)
	if !found {
		logger.Error("Could not find named cache store", zap.String("name", name))
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		logger.Error("GetCache expects a pointer to store the result")
	}
	return c.Get(ctx, key, value)
}

// SetCache initializes a value in the cache store
func SetCache(ctx context.Context, name string, key string, value interface{}) error {
	c, found := cacheWrapper[name]
	if !found {
		logger := ContextLogger(ctx)
		logger.Error("Could not find named cache store", zap.String("name", name))
	}

	return c.Set(ctx, key, value)
}

// Get retrieves a value from the cache store
func (l LocalCache) Get(ctx context.Context, key string, value interface{}) bool {
	res, found := l.Client.Get(key)
	if found {
		copier.Copy(value, res)
	}
	return found
}

// Get retrieves a value from the cache store
func (r RedisCache) Get(ctx context.Context, key string, value interface{}) bool {
	res, err := r.Client.Get(ctx, r.Prefix+key).Result()
	if err != nil {
		logger := ContextLogger(ctx)
		logger.Error("Error retrieving value from redis", zap.Error(err))
	} else {
		err = json.Unmarshal([]byte(res), &value)
	}

	return err != redis.Nil
}

// Set initializes a value in the cache store
func (l LocalCache) Set(ctx context.Context, key string, value interface{}) error {
	l.Client.Set(key, value, l.Duration)
	return nil
}

// Set initializes a value in the cache store
func (r RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	val, _ := json.Marshal(value)
	return r.Client.Set(ctx, r.Prefix+key, val, r.Duration).Err()
}

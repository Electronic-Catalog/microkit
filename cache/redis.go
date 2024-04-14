package cache

import (
	"context"
	"fmt"
	"github.com/Electronic-Catalog/microkit/logger"
	"github.com/Electronic-Catalog/microkit/logger/zap"
	"github.com/Electronic-Catalog/microkit/metric"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisCache struct {
	client               *redis.Client
	metric               metric.Metric
	reqResLogger         logger.Logger
	singleInstanceOption *redis.Options
	failOverOption       *redis.FailoverOptions
}

// NewRedisCache
// in order to use single redis, just give single address in redis addresses and for enabling redis
// sentinel mode, give cluster master-name and list of cluster members addresses
func NewRedisCache(options ...Option) (Cache, error) {
	repo := redisCache{
		client:               nil,
		metric:               metric.NewNop(),
		reqResLogger:         zap.NopLogger,
		singleInstanceOption: nil,
		failOverOption:       nil,
	}

	for _, op := range options { // configure redis cache with options
		err := op(&repo)
		if err != nil {
			return nil, err
		}
	}

	if repo.singleInstanceOption != nil {
		// set defaults
		if repo.singleInstanceOption.DialTimeout == 0 {
			repo.singleInstanceOption.DialTimeout = time.Second * 10
		}
		if repo.singleInstanceOption.MaxRetries == 0 {
			repo.singleInstanceOption.MaxRetries = 2
		}

		repo.client = redis.NewClient(repo.singleInstanceOption)
	} else if repo.failOverOption != nil {
		// set defaults
		if repo.failOverOption.DialTimeout == 0 {
			repo.failOverOption.DialTimeout = time.Second * 10
		}
		if repo.failOverOption.MaxRetries == 0 {
			repo.failOverOption.MaxRetries = 2
		}

		repo.client = redis.NewFailoverClient(repo.failOverOption)
	} else {
		return nil, fmt.Errorf("unspecified redis connection addresses - you have to use WithAddresses() option or WithConnectionStringOption() option")
	}

	pingCtx, cf := context.WithTimeout(context.Background(), time.Second*10)
	defer cf()

	err := repo.Ping(pingCtx)
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func (r *redisCache) GetKey(ctx context.Context, method string, key string) (string, error) {
	r.metric.IncrementTotal("redis", method)
	defer func(startTime time.Time) {
		r.metric.ObserveResponseTime(time.Since(startTime), "redis", method)
	}(time.Now())

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", NotFoundError
	} else if err != nil {
		r.metric.IncrementError("redis", method, err.Error())
		return "", err
	}

	return val, nil
}

func (r *redisCache) Set(ctx context.Context, method string, key string, val string, expiration time.Duration) error {
	r.metric.IncrementTotal("redis", method)
	defer func(startTime time.Time) {
		r.metric.ObserveResponseTime(time.Since(startTime), "redis", method)
	}(time.Now())

	err := r.client.Set(ctx, key, val, expiration).Err()
	if err != nil {
		r.metric.IncrementError("redis", method, err.Error())
		return err
	}

	return nil
}

func (r *redisCache) Ping(ctx context.Context) error {
	// because the ICMP ping packet are filtered in many infrastructures we are using
	// get and set command for redis to ensuring connection availability
	_, err := r.client.Get(ctx, "test").Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func (r *redisCache) RemoveKey(ctx context.Context, method string, key string) error {
	r.metric.IncrementTotal("redis", method)
	defer func(startTime time.Time) {
		r.metric.ObserveResponseTime(time.Since(startTime), "redis", method)
	}(time.Now())

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.metric.IncrementError("redis", method, err.Error())
		return err
	}

	return nil
}

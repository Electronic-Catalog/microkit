package cache

import (
	"fmt"
	"github.com/Electronic-Catalog/microkit/logger"
	"github.com/Electronic-Catalog/microkit/metric"
	"github.com/go-redis/redis/v8"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Option func(cache Cache) error

// WithMetricOption
// this option allows you to define metric for your cache instance
func WithMetricOption(metric metric.Metric) Option {
	return func(cache Cache) error {
		switch cache.(type) {
		case *redisCache:
			rd, _ := cache.(*redisCache)
			rd.metric = metric
		case *memCache:
			mc, _ := cache.(*memCache)
			mc.metric = metric
		}
		return nil
	}
}

// WithAddresses
// in order to use single instance give master name empty string, and redis addresses ad single item
func WithAddresses(masterName *string, redisAddress ...string) Option {
	return func(cache Cache) error {
		switch cache.(type) {
		case *redisCache:
			rd, _ := cache.(*redisCache)
			if len(redisAddress) == 1 {
				redisOption := redis.Options{
					Addr: redisAddress[0],
				}
				rd.singleInstanceOption = &redisOption
			} else if len(redisAddress) > 1 {
				redisOption := redis.FailoverOptions{
					MasterName:    *masterName,
					SentinelAddrs: redisAddress,
				}
				rd.failOverOption = &redisOption
			}
		}

		return nil
	}
}

func WithDbNumber(dbNumber int) Option {
	return func(cache Cache) error {
		switch cache.(type) {
		case *redisCache:
			rd, _ := cache.(*redisCache)
			if rd.failOverOption != nil {
				rd.failOverOption.DB = dbNumber
			} else if rd.singleInstanceOption != nil {
				rd.singleInstanceOption.DB = dbNumber
			}
		}

		return nil
	}
}

func WithMaxRetry(maxRetry int) Option {
	return func(cache Cache) error {
		switch cache.(type) {
		case *redisCache:
			rd, _ := cache.(*redisCache)
			if rd.failOverOption != nil {
				rd.failOverOption.MaxRetries = maxRetry
			} else if rd.singleInstanceOption != nil {
				rd.singleInstanceOption.MaxRetries = maxRetry
			}
		}

		return nil
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(cache Cache) error {
		switch cache.(type) {
		case *redisCache:
			rd, _ := cache.(*redisCache)
			if rd.failOverOption != nil {
				rd.failOverOption.DialTimeout = timeout
			} else if rd.singleInstanceOption != nil {
				rd.singleInstanceOption.DialTimeout = timeout
			}
		}

		return nil
	}
}

// WithConnectionStringOption
// this model will configure redis in the way which we have a connection string and
// we have to pars this string and create redis options
func WithConnectionStringOption(connectionString string) Option {
	return func(cache Cache) error {
		switch cache.(type) {
		case *redisCache:
			rd, _ := cache.(*redisCache)
			u, err := url.Parse(connectionString)
			if err != nil {
				return fmt.Errorf("got error %v on configuring redis with connection string", err)
			}

			db := strings.Trim(u.Path, "/")
			idb, _ := strconv.Atoi(db)
			password := ""

			if pas, ex := u.User.Password(); ex {
				password = pas
			}

			option := redis.Options{
				Addr:     u.Host,
				Password: password,
				Username: u.User.Username(),
				DB:       idb,
			}

			rd.singleInstanceOption = &option
		}

		return nil
	}
}

// WithLoggerOption
// this option enables error logging on cache instance for you
func WithLoggerOption(logger logger.Logger) Option {
	return func(cache Cache) error {
		switch cache.(type) {
		case *redisCache:
			rd, _ := cache.(*redisCache)
			rd.reqResLogger = logger
		}

		return nil
	}
}

package cache

import (
	"context"
	"github.com/Electronic-Catalog/microkit/metric"
	"sync"
	"time"
)

type item struct {
	value      string // we didn't concern us with process of encoding and decoding value
	expiration time.Time
}

type memCache struct {
	store  map[string]item
	lock   *sync.RWMutex
	metric metric.Metric
	timer  *time.Timer
}

func (m *memCache) RemoveKey(ctx context.Context, method string, key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	start := time.Now()
	defer func(start time.Time) {
		m.metric.ObserveResponseTime(time.Now().Sub(start), "mem", "del", "method", method)
	}(start)

	delete(m.store, key)

	return nil
}

func NewInMemoryCache(evictionInterval time.Duration, options ...Option) (Cache, error) {
	mm := memCache{
		store:  make(map[string]item),
		lock:   &sync.RWMutex{},
		metric: metric.NewNop(),
		timer:  time.NewTimer(evictionInterval),
	}

	for _, op := range options {
		err := op(&mm)
		if err != nil {
			return nil, err
		}
	}

	// background process of removing expired items
	go mm.evictionProcess()

	return &mm, nil
}

func (m *memCache) Ping(ctx context.Context) error {
	// NOP
	return nil
}

func (m *memCache) GetKey(ctx context.Context, method string, key string) (string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	start := time.Now()
	defer func(start time.Time) {
		m.metric.ObserveResponseTime(time.Now().Sub(start), "mem", "get", "method", method)
	}(start)
	if item, ok := m.store[key]; ok {
		return item.value, nil
	} else {
		return "", NotFoundError
	}
}

func (m *memCache) Set(ctx context.Context, method string, key string, val string, expiration time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	start := time.Now()
	defer func(start time.Time) {
		m.metric.ObserveResponseTime(time.Now().Sub(start), "mem", "set", "method", method)
	}(start)

	m.store[key] = item{
		value:      val,
		expiration: time.Now().Add(expiration),
	}

	return nil
}

func (m *memCache) evictionProcess() {
	// this implementation is so ruth and read cache in each interval and if some item expiration
	// exceeded remove them from mem cache

	for {
		<-m.timer.C
		now := time.Now()
		m.lock.Lock()
		// m.lock.RLock() // read lock
		for key, val := range m.store {
			if val.expiration.Before(now) {
				// m.lock.RUnlock() // read unlock
				// m.lock.Lock() // write lock
				delete(m.store, key)
				// m.lock.Unlock() // write unlock
				// m.lock.RLock() // read lock
			}
		}
		// m.lock.RUnlock() // read unlock
		m.lock.Unlock()
	}
}

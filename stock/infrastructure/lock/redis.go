package lock

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type RedisLocker struct {
	rdb            *redis.Client
	expiration     time.Duration
	retryNumber    int
	retryDelay     time.Duration
	watchdogStopCh map[string]chan struct{}
	mutex          sync.Mutex
}

func NewRedisLocker(host, port string, expiration time.Duration, retryNumber int, retryDelay time.Duration) *RedisLocker {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
	return &RedisLocker{
		rdb:            rdb,
		expiration:     expiration,
		retryNumber:    retryNumber,
		retryDelay:     retryDelay,
		watchdogStopCh: make(map[string]chan struct{}),
	}
}

func (r *RedisLocker) Lock(ctx context.Context, key, value string) (bool, error) {
	if value == "" {
		return false, errors.New("empty lock value")
	}

	for i := 0; i < r.retryNumber; i++ {
		ok, err := r.rdb.SetNX(ctx, key, value, r.expiration).Result()
		if err != nil {
			return false, err
		}
		if ok {
			r.mutex.Lock()
			stopCh := make(chan struct{})
			r.watchdogStopCh[key] = stopCh
			r.mutex.Unlock()

			go r.startWatchdog(ctx, key, value, stopCh)
			return true, nil
		}

		logrus.Warnf("failed to acquire lock for key %s, retrying in %v (attempt %d/%d)", key, r.retryDelay, i+1, r.retryNumber)
		if i < r.retryNumber-1 {
			time.Sleep(r.retryDelay)
		}
	}

	return false, nil
}

const unlockScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end`

func (r *RedisLocker) Unlock(ctx context.Context, key, value string) (bool, error) {
	r.mutex.Lock()
	stopCh, exists := r.watchdogStopCh[key]
	if exists {
		close(stopCh)
		delete(r.watchdogStopCh, key)
	}
	r.mutex.Unlock()

	res, err := r.rdb.Eval(ctx, unlockScript, []string{key}, value).Int64()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

// startWatchdog 每隔 expiration / 2 的时间为锁自动续期，确保在操作过程中锁不会自动释放
// 同时又保证了如果服务宕机时没有释放锁，Redis 会自动释放锁
func (r *RedisLocker) startWatchdog(ctx context.Context, key, value string, stopCh chan struct{}) {
	ticker := time.NewTicker(r.expiration / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ok, _ := r.rdb.Expire(ctx, key, r.expiration).Result(); !ok {
				return
			}
		case <-stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

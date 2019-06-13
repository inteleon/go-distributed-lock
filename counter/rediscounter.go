package counter

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type RedisCounter struct {
	key           string
	mutex         sync.Mutex
	client        *redis.Client
	expirySeconds int
}

func (dl *RedisCounter) Decr() {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()

	res := dl.client.Get(dl.key)
	_, err := res.Result()
	if err == redis.Nil {
		logrus.Warnf("Tried to decrement a non-existing (expired?) lock. This indicates that the lock expired before some processing completed.")
		return
	}

	val, _ := res.Int64()
	if val > 0 {
		dl.client.Decr(dl.key)
		logrus.Debugf("Distributed lock decremented to %v", val-1)
	} else {
		logrus.Warnf("Distributed lock was decremented with 0 or less size. This is BAD!")
	}
}

func (dl *RedisCounter) IsLocked() bool {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()
	res := dl.client.Get(dl.key)
	_, err := res.Result()
	if err == redis.Nil {
		return false
	}
	val, _ := res.Int64()
	return val > 0
}

func (dl *RedisCounter) Get() int64 {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()
	res := dl.client.Get(dl.key)
	_, err := res.Result()
	if err == redis.Nil {
		return 0
	}
	val, _ := res.Int64()
	return val
}

// Set recreates the debbie_lock with the supplied value, configured with the configured expiry in seconds.
func (dl *RedisCounter) Set(cnt int64) error {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()
	res := dl.client.Set(dl.key, cnt, time.Second*time.Duration(dl.expirySeconds))
	if res.Err() != nil {
		logrus.Errorf("Distributed lock set failed: %v", res.Err())
		return errors.Wrap(res.Err(), "Unable to set distributed lock")
	}
	logrus.Infof("Distributed lock %v set to %v", dl.key, cnt)
	return nil
}

func (dl *RedisCounter) Close() {
	dl.client.Close()
}

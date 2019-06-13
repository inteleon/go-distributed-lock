package counter

import (
	"sync"
	"time"
)

// InMemoryCounter is typically used for local testing without a proper KV store available.
type InMemoryCounter struct {
	key       string
	mutex     sync.Mutex
	count     int64
	expiresAt time.Time
}

func NewInMemoryCounter(key string, expiresAt time.Time) *InMemoryCounter {
	return &InMemoryCounter{key: key, mutex: sync.Mutex{}, count: int64(0), expiresAt: expiresAt}
}

// Decr decreases the counter by one.
func (dl *InMemoryCounter) Decr() {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()
	if dl.count > 0 {
		dl.count--
	}
}

// IsLocked returns true if the counter has a value greater than zero.
func (dl *InMemoryCounter) IsLocked() bool {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()
	return dl.count > 0
}

// Get returns the current value of the counter.
func (dl *InMemoryCounter) Get() int64 {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()
	return dl.count
}

// Set recreates the debbie_lock with the supplied value, configured with the configured expiry in seconds.
func (dl *InMemoryCounter) Set(cnt int64) error {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()
	dl.count = cnt
	return nil
}

func (dl *InMemoryCounter) Close() {
	// Noop
}

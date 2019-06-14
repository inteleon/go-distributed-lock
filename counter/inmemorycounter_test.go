package counter

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestRespectsLock(t *testing.T) {
	counter := NewInMemoryCounter("my-key", time.Now().Add(3*time.Second))

	_ = counter.Set(1)
	assert.True(t, counter.IsLocked())

	assert.EqualValues(t, 1, counter.Get())
	counter.Decr()

	assert.False(t, counter.IsLocked())
	assert.EqualValues(t, 0, counter.Get())
}

func TestDecrInDifferentGoroutines(t *testing.T) {
	counter := NewInMemoryCounter("my-key", time.Now().Add(3*time.Second))
	assert.False(t, counter.IsLocked())
	_ = counter.Set(1337)
	assert.True(t, counter.IsLocked())

	wg := sync.WaitGroup{}
	wg.Add(2)

	var c1, c2 int

	go func() {
		for counter.Get() > 0 {
			counter.Decr()
			c1++
		}
		wg.Done()
	}()
	go func() {

		for counter.Get() > 0 {
			counter.Decr()
			c2++
		}
		wg.Done()
	}()
	wg.Wait()
	assert.False(t, counter.IsLocked())
	assert.Equal(t, 1337, c1+c2)
}

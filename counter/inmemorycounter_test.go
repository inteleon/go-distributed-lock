package counter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRespectsLock(t *testing.T) {
	counter := NewInMemoryCounter("my-key", time.Now().Add(3*time.Second))

	counter.Set(1)
	assert.True(t, counter.IsLocked())

	assert.EqualValues(t, 1, counter.Get())
	counter.Decr()

	assert.False(t, counter.IsLocked())
	assert.EqualValues(t, 0, counter.Get())
}

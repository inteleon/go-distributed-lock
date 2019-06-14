package counter

import (
	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	"github.com/inteleon/go-distributed-lock/mock_redis"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"unsafe"
)

func TestRedisCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mr := mock_redis.NewMockUniversalClient(ctrl)
	mr.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(redis.NewStatusCmd())
	mr.EXPECT().Decr(gomock.Any()).Times(1)

	// Set up a sequence of expected calls to Get for this testcase
	call1 := mr.EXPECT().Get(gomock.Any()).Return(buildResp(0))
	call2 := mr.EXPECT().Get(gomock.Any()).Return(buildResp(1))
	call3 := mr.EXPECT().Get(gomock.Any()).Return(buildResp(1))
	call4 := mr.EXPECT().Get(gomock.Any()).Return(buildResp(0))
	gomock.InOrder(call1, call2, call3, call4)

	cnt := RedisCounter{
		key:           "key",
		mutex:         sync.Mutex{},
		client:        mr,
		expirySeconds: 3,
	}

	assert.False(t, cnt.IsLocked())
	err := cnt.Set(int64(1))
	assert.NoError(t, err)
	assert.True(t, cnt.IsLocked())

	cnt.Decr()
	assert.False(t, cnt.IsLocked())
}

func buildResp(i int64) *redis.StringCmd {
	stringCmd := redis.NewStringCmd()

	str := strconv.Itoa(int(i))

	// This horrible code uses reflection and unsafe tricks to mutate the unexported "val" of the StringCmd
	// so we can mock the response with our supplied value.
	reflectedStruct := reflect.ValueOf(stringCmd).Elem() // Get the elem of the return struct
	fieldToSet := reflectedStruct.Field(1)               // Get the second field of the struct (the val)
	valueToSet := reflect.ValueOf(&str).Elem()           // Get value to set as reflect.Value

	// This is h4xx to make it possible to mutate the unexported field
	fieldToSet = reflect.NewAt(fieldToSet.Type(), unsafe.Pointer(fieldToSet.UnsafeAddr())).Elem()
	fieldToSet.Set(valueToSet)
	return stringCmd
}

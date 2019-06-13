package counter

// Counter defines a lock implemented as a counter.
type Counter interface {
	IsLocked() bool
	Decr()
	Set(count int64)
	Get() int64
	Close()
}

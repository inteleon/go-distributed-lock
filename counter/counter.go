package counter

// Counter defines a lock implemented as a counter.
type Counter interface {

	// IsLocked should return the semantic meaning of the underlying locking mechanism.
	IsLocked() bool

	// Decr decrements an underlying counter. Implementations should make sure that an underlying counter never is < 0
	Decr()

	// Set sets the underlying counter to the specified value. Implementations may choose to set an expiry here.
	Set(count int64) error

	// Get returns the current counter value. A value < 0 indicates an expired or non existent lock.
	Get() int64

	// Close can be used for implementations that may need to close resources before a shutdown.
	Close()
}

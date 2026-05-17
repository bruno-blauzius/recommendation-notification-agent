package output

import "time"

// IdempotencyRepository defines the output port (driven port) for idempotency control.
// It is used to track in-flight records and prevent duplicate processing.
type IdempotencyRepository interface {
	// Set stores a key with an associated value and a TTL.
	// If the key already exists it is overwritten.
	Set(key string, value string, ttl time.Duration) error

	// Get retrieves the value associated with key.
	// Returns ("", nil) when the key does not exist.
	Get(key string) (string, error)

	// Exists reports whether the given key is present in the store.
	Exists(key string) (bool, error)
}

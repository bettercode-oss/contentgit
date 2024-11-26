package cache

import (
	"time"
)

type Cache interface {
	Set(key string, value any, expiration time.Duration)
	Get(key string) (any, bool)
}

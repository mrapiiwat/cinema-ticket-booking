package ports

import (
	"context"
	"time"
)

type LockRepositoryPort interface {
	AcquireLock(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string, value string) error
	GetLockValue(ctx context.Context, key string) (string, error)
}

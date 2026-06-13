package repositories

import (
	"context"
	"time"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/database"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
	"github.com/redis/go-redis/v9"
)

type lockRepository struct {
	client *redis.Client
}

func NewLockRepository() ports.LockRepositoryPort {
	return &lockRepository{
		client: database.RedisClient,
	}
}

func (r *lockRepository) AcquireLock(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	// It will only set the key if it does not already exist
	// Returns true if the key was set, false if it already existed
	ok, err := r.client.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (r *lockRepository) ReleaseLock(ctx context.Context, key string, value string) error {
	// Use Lua script to atomically check if the lock belongs to us before deleting it
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)

	_, err := script.Run(ctx, r.client, []string{key}, value).Result()
	return err
}

func (r *lockRepository) GetLockValue(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

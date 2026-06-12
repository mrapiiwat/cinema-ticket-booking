package repositories

import (
	"context"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/database"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
	"github.com/redis/go-redis/v9"
)

type mqRepository struct {
	client *redis.Client
}

func NewMQRepository() ports.MQRepositoryPort {
	return &mqRepository{
		client: database.RedisClient,
	}
}

func (r *mqRepository) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.client.Publish(ctx, channel, message).Err()
}

func (r *mqRepository) Subscribe(ctx context.Context, channel string, handler func(payload string)) error {
	pubsub := r.client.Subscribe(ctx, channel)

	// Start a goroutine to listen for messages
	go func() {
		defer pubsub.Close()
		for msg := range pubsub.Channel() {
			handler(msg.Payload)
		}
	}()

	return nil
}

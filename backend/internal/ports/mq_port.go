package ports

import (
	"context"
)

type MQRepositoryPort interface {
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channel string, handler func(payload string)) error
}

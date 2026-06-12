package ports

import (
	"context"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
)

type AuditRepositoryPort interface {
	Create(ctx context.Context, log *model.AuditLog) error
	FindAll(ctx context.Context) ([]model.AuditLog, error)
}

type AuditServicePort interface {
	LogEvent(ctx context.Context, eventType, details, userID string) (*model.AuditLog, error)
	GetLogs(ctx context.Context) ([]model.AuditLog, error)
}

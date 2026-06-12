package services

import (
	"context"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
)

type auditService struct {
	auditRepo ports.AuditRepositoryPort
}

func NewAuditService(auditRepo ports.AuditRepositoryPort) ports.AuditServicePort {
	return &auditService{
		auditRepo: auditRepo,
	}
}

func (s *auditService) LogEvent(ctx context.Context, eventType, details, userID string) (*model.AuditLog, error) {
	log := &model.AuditLog{
		EventType: eventType,
		Details:   details,
		UserID:    userID,
	}
	if err := s.auditRepo.Create(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *auditService) GetLogs(ctx context.Context) ([]model.AuditLog, error) {
	return s.auditRepo.FindAll(ctx)
}

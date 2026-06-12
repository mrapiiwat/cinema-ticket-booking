package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
)

func StartAuditLogConsumer(mqRepo ports.MQRepositoryPort, auditService ports.AuditServicePort, wsHub *WSHub) {
	err := mqRepo.Subscribe(context.Background(), "audit_logs", func(payload string) {
		var logMsg model.AuditLog
		if err := json.Unmarshal([]byte(payload), &logMsg); err != nil {
			log.Printf("Failed to unmarshal audit log message: %v\n", err)
			return
		}

		savedLog, err := auditService.LogEvent(context.Background(), logMsg.EventType, logMsg.Details, logMsg.UserID)
		if err != nil {
			log.Printf("Failed to save audit log: %v\n", err)
			return
		}
		if wsHub != nil {
			wsHub.BroadcastMessage(WSMessage{
				Type:    "AUDIT_LOG_CREATED",
				Payload: savedLog,
			})
		}
	})

	if err != nil {
		log.Printf("Failed to start audit log consumer: %v\n", err)
	}
}

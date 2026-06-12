package repositories

import (
	"context"
	"time"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/database"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type auditRepository struct {
	collection *mongo.Collection
}

func NewAuditRepository() ports.AuditRepositoryPort {
	return &auditRepository{
		collection: database.GetCollection("audit_logs"),
	}
}

func (r *auditRepository) Create(ctx context.Context, log *model.AuditLog) error {
	log.Timestamp = time.Now()
	res, err := r.collection.InsertOne(ctx, log)
	if err == nil {
		log.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

func (r *auditRepository) FindAll(ctx context.Context) ([]model.AuditLog, error) {
	var logs []model.AuditLog

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

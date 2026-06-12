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

type bookingRepository struct {
	collection *mongo.Collection
}

func NewBookingRepository() ports.BookingRepositoryPort {
	return &bookingRepository{
		collection: database.GetCollection("bookings"),
	}
}

func (r *bookingRepository) Create(ctx context.Context, booking *model.Booking) error {
	now := time.Now()
	booking.CreatedAt = now
	booking.UpdatedAt = now
	if booking.Status == "" {
		booking.Status = model.BookingStatusPending
	}
	res, err := r.collection.InsertOne(ctx, booking)
	if err == nil {
		booking.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

func (r *bookingRepository) FindByID(ctx context.Context, id string) (*model.Booking, error) {
	var booking model.Booking
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	set := bson.M{"status": status, "updated_at": time.Now()}
	if status == model.BookingStatusSuccess {
		now := time.Now()
		set["confirmed_at"] = now
	}
	update := bson.M{"$set": set}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *bookingRepository) FindByUserID(ctx context.Context, userID string) ([]model.Booking, error) {
	var bookings []model.Booking
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) FindAll(ctx context.Context, filter model.BookingFilter) ([]model.Booking, error) {
	query := bson.M{}
	if filter.UserID != "" {
		query["user_id"] = filter.UserID
	}
	if filter.ShowtimeID != "" {
		objID, err := primitive.ObjectIDFromHex(filter.ShowtimeID)
		if err != nil {
			return nil, err
		}
		query["showtime_id"] = objID
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []model.Booking
	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) FindExpiredPending(ctx context.Context, now time.Time) ([]model.Booking, error) {
	query := bson.M{
		"status":          model.BookingStatusPending,
		"lock_expires_at": bson.M{"$lte": now},
	}
	opts := options.Find().SetSort(bson.D{{Key: "lock_expires_at", Value: 1}})
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []model.Booking
	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

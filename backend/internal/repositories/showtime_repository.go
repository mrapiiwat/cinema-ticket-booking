package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/database"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrSeatUnavailable = errors.New("seat is not available")

type showtimeRepository struct {
	collection *mongo.Collection
}

func NewShowtimeRepository() ports.ShowtimeRepositoryPort {
	return &showtimeRepository{
		collection: database.GetCollection("showtimes"),
	}
}

func (r *showtimeRepository) FindByMovieID(ctx context.Context, movieID string) ([]model.Showtime, error) {
	var showtimes []model.Showtime
	objID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"movie_id": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &showtimes); err != nil {
		return nil, err
	}
	return showtimes, nil
}

func (r *showtimeRepository) FindByID(ctx context.Context, id string) (*model.Showtime, error) {
	var showtime model.Showtime
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&showtime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &showtime, nil
}

func (r *showtimeRepository) Create(ctx context.Context, showtime *model.Showtime) error {
	res, err := r.collection.InsertOne(ctx, showtime)
	if err == nil {
		showtime.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

func (r *showtimeRepository) UpdateSeatStatus(ctx context.Context, showtimeID string, seatNo string, status string, lockedBy *string, lockedUntil *time.Time) error {
	objID, err := primitive.ObjectIDFromHex(showtimeID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID, "seats.seat_no": seatNo}
	set := bson.M{"seats.$.status": status}
	update := bson.M{"$set": set}
	if lockedBy != nil {
		set["seats.$.locked_by"] = *lockedBy
	} else {
		update["$unset"] = bson.M{"seats.$.locked_by": ""}
	}
	if lockedUntil != nil {
		set["seats.$.locked_until"] = *lockedUntil
	} else {
		unset, ok := update["$unset"].(bson.M)
		if !ok {
			unset = bson.M{}
			update["$unset"] = unset
		}
		unset["seats.$.locked_until"] = ""
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *showtimeRepository) LockSeatIfAvailable(ctx context.Context, showtimeID string, seatNo string, userID string, lockedUntil time.Time) error {
	objID, err := primitive.ObjectIDFromHex(showtimeID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
		"seats": bson.M{"$elemMatch": bson.M{
			"seat_no": seatNo,
			"$or": []bson.M{
				{"status": model.SeatStatusAvailable},
				{"status": model.SeatStatusLocked, "locked_until": bson.M{"$lte": time.Now()}},
			},
		}},
	}
	update := bson.M{"$set": bson.M{
		"seats.$.status":       model.SeatStatusLocked,
		"seats.$.locked_by":    userID,
		"seats.$.locked_until": lockedUntil,
	}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrSeatUnavailable
	}
	return nil
}

func (r *showtimeRepository) ReleaseSeatLock(ctx context.Context, showtimeID string, seatNo string, userID string) error {
	objID, err := primitive.ObjectIDFromHex(showtimeID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":             objID,
		"seats.seat_no":   seatNo,
		"seats.status":    model.SeatStatusLocked,
		"seats.locked_by": userID,
	}
	update := bson.M{
		"$set": bson.M{"seats.$.status": model.SeatStatusAvailable},
		"$unset": bson.M{
			"seats.$.locked_by":    "",
			"seats.$.locked_until": "",
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *showtimeRepository) ReleaseExpiredLocks(ctx context.Context, showtimeID string, now time.Time) error {
	objID, err := primitive.ObjectIDFromHex(showtimeID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{"seats.$[expired].status": model.SeatStatusAvailable},
		"$unset": bson.M{
			"seats.$[expired].locked_by":    "",
			"seats.$[expired].locked_until": "",
		},
	}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{
				"expired.status":       model.SeatStatusLocked,
				"expired.locked_until": bson.M{"$lte": now},
			},
		},
	})

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update, opts)
	return err
}

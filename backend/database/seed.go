package database

import (
	"context"
	"time"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type seedLogger interface {
	Infof(format string, args ...interface{})
}

func SeedDefaults(logger seedLogger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ensureIndexes(ctx); err != nil {
		return err
	}

	count, err := GetCollection("movies").CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	movies := []model.Movie{
		{Title: "Dune: Part Two", DurationMins: 166, PosterURL: "https://i.redd.it/99a79xg2yfxa1.jpg", CreatedAt: time.Now()},
		{Title: "Oppenheimer", DurationMins: 180, PosterURL: "https://upload.wikimedia.org/wikipedia/donate/5/5b/Oppenheimer_poster.jpg?utm_source=donate.wikimedia.org&utm_campaign=index&utm_content=original", CreatedAt: time.Now()},
		{Title: "Spider-Man: No Way Home", DurationMins: 148, PosterURL: "https://m.media-amazon.com/images/S/pv-target-images/0a76c4a1044229d184bb09480b314150340ee8c45de8365c062a6e0f20201614.jpg", CreatedAt: time.Now()},
	}

	movieDocs := make([]interface{}, 0, len(movies))
	for i := range movies {
		movieDocs = append(movieDocs, movies[i])
	}
	res, err := GetCollection("movies").InsertMany(ctx, movieDocs)
	if err != nil {
		return err
	}
	for i, id := range res.InsertedIDs {
		movies[i].ID = id.(primitive.ObjectID)
	}

	now := time.Now()
	showtimes := []interface{}{
		model.Showtime{MovieID: movies[0].ID, StartTime: now.Add(2 * time.Hour), PricePerSeat: 220, Seats: defaultSeats()},
		model.Showtime{MovieID: movies[0].ID, StartTime: now.Add(6 * time.Hour), PricePerSeat: 240, Seats: defaultSeats()},
		model.Showtime{MovieID: movies[1].ID, StartTime: now.Add(3 * time.Hour), PricePerSeat: 230, Seats: defaultSeats()},
		model.Showtime{MovieID: movies[2].ID, StartTime: now.Add(24 * time.Hour), PricePerSeat: 190, Seats: defaultSeats()},
	}

	if _, err := GetCollection("showtimes").InsertMany(ctx, showtimes); err != nil {
		return err
	}

	if logger != nil {
		logger.Infof("Seeded %d movies and %d showtimes", len(movies), len(showtimes))
	}
	return nil
}

func ensureIndexes(ctx context.Context) error {
	if _, err := GetCollection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "google_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return err
	}

	if _, err := GetCollection("bookings").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}, {Key: "lock_expires_at", Value: 1}},
	}); err != nil {
		return err
	}

	_, err := GetCollection("audit_logs").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "timestamp", Value: -1}},
	})
	return err
}

func defaultSeats() []model.Seat {
	rows := []string{"A", "B", "C", "D", "E", "F"}
	seats := make([]model.Seat, 0, len(rows)*8)
	for _, row := range rows {
		for number := 1; number <= 8; number++ {
			seats = append(seats, model.Seat{
				SeatNo: row + string(rune('0'+number)),
				Status: model.SeatStatusAvailable,
			})
		}
	}
	return seats
}

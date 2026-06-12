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
)

type movieRepository struct {
	collection *mongo.Collection
}

func NewMovieRepository() ports.MovieRepositoryPort {
	return &movieRepository{
		collection: database.GetCollection("movies"),
	}
}

func (r *movieRepository) FindAll(ctx context.Context) ([]model.Movie, error) {
	var movies []model.Movie
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *movieRepository) FindByID(ctx context.Context, id string) (*model.Movie, error) {
	var movie model.Movie
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &movie, nil
}

func (r *movieRepository) Create(ctx context.Context, movie *model.Movie) error {
	movie.CreatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, movie)
	if err == nil {
		movie.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

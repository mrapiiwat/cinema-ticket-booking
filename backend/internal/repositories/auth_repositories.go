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

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository() ports.UserRepositoryPort {
	return &userRepository{
		collection: database.GetCollection("users"),
	}
}

func (r *userRepository) FindByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	var user model.User
	filter := bson.M{"google_id": googleID}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	user.CreatedAt = time.Now()
	if user.Role == "" {
		user.Role = model.RoleUser
	}

	res, err := r.collection.InsertOne(ctx, user)
	if err == nil {
		user.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

func (r *userRepository) UpdateGoogleProfile(ctx context.Context, user *model.User) error {
	if user == nil || user.GoogleID == "" {
		return nil
	}

	update := bson.M{"$set": bson.M{
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}}
	_, err := r.collection.UpdateOne(ctx, bson.M{"google_id": user.GoogleID}, update)
	return err
}

package repository

import (
	"context"
	"ddd-user-service/internal/domain"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

type mongoUser struct {
	ID       string `bson:"_id"`
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Username string `bson:"username"`
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	collection := db.Collection("users")

	// Create unique indexes for email and username
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	collection.Indexes().CreateMany(context.Background(), indexModels)

	return &MongoUserRepository{
		collection: collection,
	}
}

func (r *MongoUserRepository) Save(ctx context.Context, user *domain.User) error {
	mongoUser := mongoUser{
		ID:       user.ID.String(),
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
	}

	_, err := r.collection.InsertOne(ctx, mongoUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrEmailExists
		}
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

func (r *MongoUserRepository) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	var mongoUser mongoUser
	err := r.collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&mongoUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return r.mongoUserToDomain(&mongoUser), nil
}

func (r *MongoUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var mongoUser mongoUser
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&mongoUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return r.mongoUserToDomain(&mongoUser), nil
}

func (r *MongoUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var mongoUser mongoUser
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&mongoUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return r.mongoUserToDomain(&mongoUser), nil
}

func (r *MongoUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var mongoUser mongoUser
		if err := cursor.Decode(&mongoUser); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, r.mongoUserToDomain(&mongoUser))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, user *domain.User) error {
	update := bson.M{
		"$set": bson.M{
			"name":     user.Name,
			"email":    user.Email,
			"username": user.Username,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID.String()}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrEmailExists
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *MongoUserRepository) Delete(ctx context.Context, id domain.UserID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id.String()})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *MongoUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

func (r *MongoUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return count > 0, nil
}

func (r *MongoUserRepository) mongoUserToDomain(mongoUser *mongoUser) *domain.User {
	return &domain.User{
		ID:       domain.UserID(mongoUser.ID),
		Name:     mongoUser.Name,
		Email:    mongoUser.Email,
		Username: mongoUser.Username,
	}
}

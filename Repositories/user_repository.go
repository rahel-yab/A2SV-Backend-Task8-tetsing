package Repositories

import (
	"context"
	"task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) Domain.UserRepository {
	return &mongoUserRepository{
		collection: collection,
	}
}

func (r *mongoUserRepository) AddUser(ctx context.Context, user *Domain.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *mongoUserRepository) GetUserByEmail(ctx context.Context, email string) (*Domain.User, error) {
	var user Domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) GetUserByUsername(ctx context.Context, username string) (*Domain.User, error) {
	var user Domain.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) IsUsersCollectionEmpty(ctx context.Context) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	return count == 0, err
}

func (r *mongoUserRepository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	return count > 0, err
}

func (r *mongoUserRepository) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"username": username})
	return count > 0, err
}

func (r *mongoUserRepository) PromoteUserToAdmin(ctx context.Context, identifier string) error {
	filter := bson.M{"$or": []bson.M{{"username": identifier}, {"email": identifier}}}
	update := bson.M{"$set": bson.M{"role": "admin"}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

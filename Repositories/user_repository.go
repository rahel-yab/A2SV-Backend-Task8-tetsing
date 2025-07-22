package Repositories

import (
	"context"
	"task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	AddUser(user Domain.User) error
	GetUserByEmail(email string) (*Domain.User, error)
	GetUserByUsername(username string) (*Domain.User, error)
	IsUsersCollectionEmpty() (bool, error)
	UserExistsByEmail(email string) (bool, error)
	UserExistsByUsername(username string) (bool, error)
	PromoteUserToAdmin(identifier string) error
}

type MongoUserRepository struct {
	Collection *mongo.Collection
}

func (r *MongoUserRepository) AddUser(user Domain.User) error {
	_, err := r.Collection.InsertOne(context.Background(), user)
	return err
}

func (r *MongoUserRepository) GetUserByEmail(email string) (*Domain.User, error) {
	var user Domain.User
	err := r.Collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) GetUserByUsername(username string) (*Domain.User, error) {
	var user Domain.User
	err := r.Collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) IsUsersCollectionEmpty() (bool, error) {
	count, err := r.Collection.CountDocuments(context.Background(), bson.M{})
	return count == 0, err
}

func (r *MongoUserRepository) UserExistsByEmail(email string) (bool, error) {
	count, err := r.Collection.CountDocuments(context.Background(), bson.M{"email": email})
	return count > 0, err
}

func (r *MongoUserRepository) UserExistsByUsername(username string) (bool, error) {
	count, err := r.Collection.CountDocuments(context.Background(), bson.M{"username": username})
	return count > 0, err
}

func (r *MongoUserRepository) PromoteUserToAdmin(identifier string) error {
	filter := bson.M{"$or": []bson.M{{"username": identifier}, {"email": identifier}}}
	update := bson.M{"$set": bson.M{"role": "admin"}}
	_, err := r.Collection.UpdateOne(context.Background(), filter, update)
	return err
}

package repositories

import (
	"context"
	"task_manager/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserDAO (Data Access Object) is the MongoDB representation of a user
// Used for database serialization/deserialization with bson tags
type UserDAO struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
	Role     string `bson:"role"`
}

func userToDAO(user *domain.User) *UserDAO {
	return &UserDAO{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}
}

func daoToUser(dao *UserDAO) *domain.User {
	return &domain.User{
		Username: dao.Username,
		Email:    dao.Email,
		Password: dao.Password,
		Role:     dao.Role,
	}
}

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) domain.IUserRepository {
	db := client.Database("task_manager")
	return &mongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *mongoUserRepository) AddUser(ctx context.Context, user *domain.User) error {
    dao := userToDAO(user)
    _, err := r.collection.InsertOne(ctx, dao)
    return err
}

func (r *mongoUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dao UserDAO
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&dao)
	if err != nil {
		return nil, err
	}
	return daoToUser(&dao), nil
}

func (r *mongoUserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var dao UserDAO
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&dao)
	if err != nil {
		return nil, err
	}
	return daoToUser(&dao), nil
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
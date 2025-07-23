package Repositories

import (
	"context"
	"task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) Domain.TaskRepository {
	return &mongoTaskRepository{
		collection: collection,
	}
}

func (r *mongoTaskRepository) AddTask(ctx context.Context, task *Domain.Task) error {
	_, err := r.collection.InsertOne(ctx, task)
	return err
}

func (r *mongoTaskRepository) GetAllTasks(ctx context.Context) ([]Domain.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var tasks []Domain.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *mongoTaskRepository) GetTaskByID(ctx context.Context, id string) (*Domain.Task, error) {
	var task Domain.Task
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *mongoTaskRepository) UpdateTask(ctx context.Context, task *Domain.Task) error {
	filter := bson.M{"id": task.ID}
	update := bson.M{"$set": task}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoTaskRepository) DeleteTask(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
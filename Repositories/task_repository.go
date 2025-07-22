package Repositories

import (
	"context"
	"errors"
	"task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository interface {
	GetAllTasks() ([]Domain.Task, error)
	GetTaskByID(id string) (*Domain.Task, error)
	AddTask(task Domain.Task) error
	UpdateTask(task Domain.Task) error
	DeleteTask(id string) error
}

type MongoTaskRepository struct {
	Collection *mongo.Collection
}

func (r *MongoTaskRepository) GetAllTasks() ([]Domain.Task, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var tasks []Domain.Task
	if err := cursor.All(context.Background(), &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *MongoTaskRepository) GetTaskByID(id string) (*Domain.Task, error) {
	var task Domain.Task
	err := r.Collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *MongoTaskRepository) AddTask(task Domain.Task) error {
	// Check if a task with the same id exists
	count, err := r.Collection.CountDocuments(context.Background(), bson.M{"id": task.ID})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("task with this id already exists")
	}
	_, err = r.Collection.InsertOne(context.Background(), task)
	return err
}

func (r *MongoTaskRepository) UpdateTask(task Domain.Task) error {
	filter := bson.M{"id": task.ID}
	update := bson.M{"$set": task}
	_, err := r.Collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *MongoTaskRepository) DeleteTask(id string) error {
	_, err := r.Collection.DeleteOne(context.Background(), bson.M{"id": id})
	return err
}
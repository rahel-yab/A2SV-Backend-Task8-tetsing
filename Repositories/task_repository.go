package repositories

import (
	"context"
	"task_manager/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskDAO (Data Access Object) is the MongoDB representation of a task
// Used for database serialization/deserialization with bson tags
type TaskDAO struct {
	ID          string    `bson:"_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	DueDate     time.Time `bson:"due_date"`
	Status      string    `bson:"status"`
}

func taskToDAO(task *domain.Task) *TaskDAO {
	return &TaskDAO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      task.Status,
	}
}

func daoToTask(dao *TaskDAO) *domain.Task {
	return &domain.Task{
		ID:          dao.ID,
		Title:       dao.Title,
		Description: dao.Description,
		DueDate:     dao.DueDate,
		Status:      dao.Status,
	}
}

type mongoTaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(client *mongo.Client) domain.ITaskRepository {
	db := client.Database("task_manager")
	return &mongoTaskRepository{
		collection: db.Collection("tasks"),
	}
}

func (r *mongoTaskRepository) AddTask(ctx context.Context, task *domain.Task) error {
	dao := taskToDAO(task)
	_, err := r.collection.InsertOne(ctx, dao)
	return err
}

func (r *mongoTaskRepository) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var daos []TaskDAO
	if err := cursor.All(ctx, &daos); err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, len(daos))
	for i, dao := range daos {
		tasks[i] = *daoToTask(&dao)
	}
	return tasks, nil
}

func (r *mongoTaskRepository) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	var dao TaskDAO
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&dao)
	if err != nil {
		return nil, err
	}
	return daoToTask(&dao), nil
}

func (r *mongoTaskRepository) UpdateTask(ctx context.Context, task *domain.Task) error {
	filter := bson.M{"_id": task.ID}
	update := bson.M{"$set": taskToDAO(task)}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoTaskRepository) DeleteTask(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
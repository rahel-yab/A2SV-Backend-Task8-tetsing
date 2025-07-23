package Usecases

import (
	"context"
	"task_manager/Domain"
	"time"
)

type TaskUsecase struct {
	taskRepository Domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecase(taskRepository Domain.TaskRepository, timeout time.Duration) *TaskUsecase {
	return &TaskUsecase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}

func (tu *TaskUsecase) Create(c context.Context, task *Domain.Task) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.AddTask(ctx, task)
}

func (tu *TaskUsecase) GetAllTasks(c context.Context) ([]Domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.GetAllTasks(ctx)
}

func (tu *TaskUsecase) GetTaskByID(c context.Context, id string) (*Domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.GetTaskByID(ctx, id)
}

func (tu *TaskUsecase) UpdateTask(c context.Context, task *Domain.Task) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.UpdateTask(ctx, task)
}

func (tu *TaskUsecase) DeleteTask(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.DeleteTask(ctx, id)
}



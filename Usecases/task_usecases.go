package usecases

import (
	"context"
	"errors"
	"task_manager/domain"
	"time"
)

type TaskUsecase struct {
	taskRepository domain.ITaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecase(taskRepository domain.ITaskRepository, timeout time.Duration) *TaskUsecase {
	return &TaskUsecase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}

func (tu *TaskUsecase) Create(c context.Context, task *domain.Task) error {
	// Validate task data
	if task == nil {
		return errors.New("task cannot be nil")
	}
	if task.Title == "" {
		return errors.New("title is required")
	}
	if task.Description == "" {
		return errors.New("description is required")
	}
	if task.Status == "" {
		return errors.New("status is required")
	}
	if task.DueDate.IsZero() {
		return errors.New("due date is required")
	}
	
	// Validate status values
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
	statusValid := false
	for _, status := range validStatuses {
		if task.Status == status {
			statusValid = true
			break
		}
	}
	if !statusValid {
		return errors.New("invalid status: must be pending, in_progress, completed, or cancelled")
	}

	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.AddTask(ctx, task)
}

func (tu *TaskUsecase) GetAllTasks(c context.Context) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.GetAllTasks(ctx)
}

func (tu *TaskUsecase) GetTaskByID(c context.Context, id string) (*domain.Task, error) {
	if id == "" {
		return nil, errors.New("task ID is required")
	}
	
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.GetTaskByID(ctx, id)
}

func (tu *TaskUsecase) UpdateTask(c context.Context, task *domain.Task) error {
	// Validate task data
	if task == nil {
		return errors.New("task cannot be nil")
	}
	if task.ID == "" {
		return errors.New("task ID is required")
	}
	if task.Title == "" {
		return errors.New("title is required")
	}
	if task.Description == "" {
		return errors.New("description is required")
	}
	if task.Status == "" {
		return errors.New("status is required")
	}
	if task.DueDate.IsZero() {
		return errors.New("due date is required")
	}
	
	// Validate status values
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
	statusValid := false
	for _, status := range validStatuses {
		if task.Status == status {
			statusValid = true
			break
		}
	}
	if !statusValid {
		return errors.New("invalid status: must be pending, in_progress, completed, or cancelled")
	}

	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.UpdateTask(ctx, task)
}

func (tu *TaskUsecase) DeleteTask(c context.Context, id string) error {
	if id == "" {
		return errors.New("task ID is required")
	}
	
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.DeleteTask(ctx, id)
}



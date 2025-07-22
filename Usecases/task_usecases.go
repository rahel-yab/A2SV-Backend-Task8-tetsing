package Usecases

import (
	"task_manager/Domain"
	"task_manager/Repositories"
)

type TaskUsecase struct {
	Repo Repositories.TaskRepository
}

func (u *TaskUsecase) GetAllTasks() ([]Domain.Task, error) {
	return u.Repo.GetAllTasks()
}

func (u *TaskUsecase) GetTaskByID(id string) (*Domain.Task, error) {
	return u.Repo.GetTaskByID(id)
}

func (u *TaskUsecase) AddTask(task Domain.Task) error {
	return u.Repo.AddTask(task)
}

func (u *TaskUsecase) UpdateTask(task Domain.Task) error {
	return u.Repo.UpdateTask(task)
}

func (u *TaskUsecase) DeleteTask(id string) error {
	return u.Repo.DeleteTask(id)
}



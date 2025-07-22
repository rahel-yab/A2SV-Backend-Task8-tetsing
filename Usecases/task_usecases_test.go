package Usecases

import (
	"errors"
	"task_manager/Domain"
	"testing"
)

type mockTaskRepo struct {
	tasks map[string]Domain.Task
}

func (m *mockTaskRepo) GetAllTasks() ([]Domain.Task, error) {
	var out []Domain.Task
	for _, t := range m.tasks {
		out = append(out, t)
	}
	return out, nil
}
func (m *mockTaskRepo) GetTaskByID(id string) (*Domain.Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &t, nil
}
func (m *mockTaskRepo) AddTask(task Domain.Task) error {
	m.tasks[task.ID] = task
	return nil
}
func (m *mockTaskRepo) UpdateTask(task Domain.Task) error {
	if _, ok := m.tasks[task.ID]; !ok {
		return errors.New("not found")
	}
	m.tasks[task.ID] = task
	return nil
}
func (m *mockTaskRepo) DeleteTask(id string) error {
	if _, ok := m.tasks[id]; !ok {
		return errors.New("not found")
	}
	delete(m.tasks, id)
	return nil
}

func TestAddAndGetTask(t *testing.T) {
	repo := &mockTaskRepo{tasks: make(map[string]Domain.Task)}
	uc := &TaskUsecase{Repo: repo}

task := Domain.Task{ID: "1", Title: "Test Task"}
	if err := uc.AddTask(task); err != nil {
		t.Errorf("add task failed: %v", err)
	}
	got, err := uc.GetTaskByID("1")
	if err != nil || got.ID != "1" {
		t.Errorf("get task failed: %v", err)
	}
}

func TestUpdateAndDeleteTask(t *testing.T) {
	repo := &mockTaskRepo{tasks: make(map[string]Domain.Task)}
	uc := &TaskUsecase{Repo: repo}

task := Domain.Task{ID: "1", Title: "Test Task"}
	uc.AddTask(task)

task.Title = "Updated"
	if err := uc.UpdateTask(task); err != nil {
		t.Errorf("update failed: %v", err)
	}
	if err := uc.DeleteTask("1"); err != nil {
		t.Errorf("delete failed: %v", err)
	}
	if _, err := uc.GetTaskByID("1"); err == nil {
		t.Errorf("expected not found after delete")
	}
} 
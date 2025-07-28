package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"task_manager/domain"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockTaskRepository is a mock implementation of ITaskRepository
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) AddTask(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) UpdateTask(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) DeleteTask(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// TaskUsecaseTestSuite is a test suite for TaskUsecase
type TaskUsecaseTestSuite struct {
	suite.Suite
	mockRepo *MockTaskRepository
	usecase  *TaskUsecase
	ctx      context.Context
}

// SetupSuite runs once before all tests in the suite
func (suite *TaskUsecaseTestSuite) SetupSuite() {
	suite.ctx = context.Background()
}

// SetupTest runs before each test
func (suite *TaskUsecaseTestSuite) SetupTest() {
	suite.mockRepo = new(MockTaskRepository)
	suite.usecase = NewTaskUsecase(suite.mockRepo, 5*time.Second)
}

// TearDownTest runs after each test
func (suite *TaskUsecaseTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestCreateTaskSuite tests the Create method
func (suite *TaskUsecaseTestSuite) TestCreateTaskSuite() {
	suite.Run("Success", func() {
		task := &domain.Task{
			ID:          "task123",
			Title:       "Test Task",
			Description: "Test Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
		}

		suite.mockRepo.On("AddTask", mock.AnythingOfType("*context.timerCtx"), task).Return(nil)

		err := suite.usecase.Create(suite.ctx, task)

		suite.NoError(err)
	})

	suite.Run("NilTask", func() {
		err := suite.usecase.Create(suite.ctx, nil)

		suite.Error(err)
		suite.Equal("task cannot be nil", err.Error())
	})

	suite.Run("EmptyTitle", func() {
		task := &domain.Task{
			ID:          "task123",
			Title:       "", // Empty title
			Description: "Test Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
		}

		err := suite.usecase.Create(suite.ctx, task)

		suite.Error(err)
		suite.Equal("title is required", err.Error())
	})

	suite.Run("EmptyDescription", func() {
		task := &domain.Task{
			ID:          "task123",
			Title:       "Test Task",
			Description: "", // Empty description
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
		}

		err := suite.usecase.Create(suite.ctx, task)

		suite.Error(err)
		suite.Equal("description is required", err.Error())
	})

	suite.Run("EmptyStatus", func() {
		task := &domain.Task{
			ID:          "task123",
			Title:       "Test Task",
			Description: "Test Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "", // Empty status
		}

		err := suite.usecase.Create(suite.ctx, task)

		suite.Error(err)
		suite.Equal("status is required", err.Error())
	})

	suite.Run("ZeroDueDate", func() {
		task := &domain.Task{
			ID:          "task123",
			Title:       "Test Task",
			Description: "Test Description",
			DueDate:     time.Time{}, // Zero due date
			Status:      "pending",
		}

		err := suite.usecase.Create(suite.ctx, task)

		suite.Error(err)
		suite.Equal("due date is required", err.Error())
	})

	suite.Run("InvalidStatus", func() {
		task := &domain.Task{
			ID:          "task123",
			Title:       "Test Task",
			Description: "Test Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "invalid_status", // Invalid status
		}

		err := suite.usecase.Create(suite.ctx, task)

		suite.Error(err)
		suite.Equal("invalid status: must be pending, in_progress, completed, or cancelled", err.Error())
	})

	suite.Run("ValidStatuses", func() {
		validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}

		for _, status := range validStatuses {
			suite.Run("Status_"+status, func() {
				task := &domain.Task{
					ID:          "task123",
					Title:       "Test Task",
					Description: "Test Description",
					DueDate:     time.Now().Add(24 * time.Hour),
					Status:      status,
				}

				suite.mockRepo.On("AddTask", mock.AnythingOfType("*context.timerCtx"), task).Return(nil)

				err := suite.usecase.Create(suite.ctx, task)

				suite.NoError(err)
			})
		}
	})
}

// TestGetAllTasksSuite tests the GetAllTasks method
func (suite *TaskUsecaseTestSuite) TestGetAllTasksSuite() {
	suite.Run("Success", func() {
		expectedTasks := []domain.Task{
			{
				ID:          "task1",
				Title:       "Task 1",
				Description: "Description 1",
				DueDate:     time.Now().Add(24 * time.Hour),
				Status:      "pending",
			},
			{
				ID:          "task2",
				Title:       "Task 2",
				Description: "Description 2",
				DueDate:     time.Now().Add(48 * time.Hour),
				Status:      "completed",
			},
		}

		suite.mockRepo.On("GetAllTasks", mock.AnythingOfType("*context.timerCtx")).Return(expectedTasks, nil)

		tasks, err := suite.usecase.GetAllTasks(suite.ctx)

		suite.NoError(err)
		suite.Equal(expectedTasks, tasks)
		suite.Len(tasks, 2)
	})

	suite.Run("Error", func() {
		// Create a new mock for this specific test to avoid interference
		mockRepo := new(MockTaskRepository)
		usecase := NewTaskUsecase(mockRepo, 5*time.Second)
		
		expectedError := errors.New("database connection failed")
		mockRepo.On("GetAllTasks", mock.AnythingOfType("*context.timerCtx")).Return([]domain.Task{}, expectedError)

		tasks, err := usecase.GetAllTasks(suite.ctx)

		suite.Error(err)
		suite.Empty(tasks)
		suite.Equal(expectedError, err)
		mockRepo.AssertExpectations(suite.T())
	})
}

// TestGetTaskByIDSuite tests the GetTaskByID method
func (suite *TaskUsecaseTestSuite) TestGetTaskByIDSuite() {
	suite.Run("Success", func() {
		expectedTask := &domain.Task{
			ID:          "task123",
			Title:       "Test Task",
			Description: "Test Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
		}

		suite.mockRepo.On("GetTaskByID", mock.AnythingOfType("*context.timerCtx"), "task123").Return(expectedTask, nil)

		task, err := suite.usecase.GetTaskByID(suite.ctx, "task123")

		suite.NoError(err)
		suite.Equal(expectedTask, task)
	})

	suite.Run("EmptyID", func() {
		task, err := suite.usecase.GetTaskByID(suite.ctx, "")

		suite.Error(err)
		suite.Nil(task)
		suite.Equal("task ID is required", err.Error())
	})

	suite.Run("NotFound", func() {
		expectedError := errors.New("task not found")
		suite.mockRepo.On("GetTaskByID", mock.AnythingOfType("*context.timerCtx"), "nonexistent").Return(nil, expectedError)

		task, err := suite.usecase.GetTaskByID(suite.ctx, "nonexistent")

		suite.Error(err)
		suite.Nil(task)
		suite.Equal(expectedError, err)
	})
}

// TestUpdateTaskSuite tests the UpdateTask method
func (suite *TaskUsecaseTestSuite) TestUpdateTaskSuite() {
	suite.Run("Success", func() {
		task := &domain.Task{
			ID:          "task123",
			Title:       "Updated Task",
			Description: "Updated Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "in_progress",
		}

		suite.mockRepo.On("UpdateTask", mock.AnythingOfType("*context.timerCtx"), task).Return(nil)

		err := suite.usecase.UpdateTask(suite.ctx, task)

		suite.NoError(err)
	})

	suite.Run("NilTask", func() {
		err := suite.usecase.UpdateTask(suite.ctx, nil)

		suite.Error(err)
		suite.Equal("task cannot be nil", err.Error())
	})

	suite.Run("EmptyID", func() {
		task := &domain.Task{
			ID:          "", // Empty ID
			Title:       "Updated Task",
			Description: "Updated Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "in_progress",
		}

		err := suite.usecase.UpdateTask(suite.ctx, task)

		suite.Error(err)
		suite.Equal("task ID is required", err.Error())
	})

	suite.Run("EmptyTitle", func() {
		task := &domain.Task{
			ID:          "task123",
			Title:       "", // Empty title
			Description: "Updated Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "in_progress",
		}

		err := suite.usecase.UpdateTask(suite.ctx, task)

		suite.Error(err)
		suite.Equal("title is required", err.Error())
	})

	suite.Run("Error", func() {
		// Create a new mock for this specific test to avoid interference
		mockRepo := new(MockTaskRepository)
		usecase := NewTaskUsecase(mockRepo, 5*time.Second)
		
		task := &domain.Task{
			ID:          "task123",
			Title:       "Updated Task",
			Description: "Updated Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "in_progress",
		}

		expectedError := errors.New("task not found")
		mockRepo.On("UpdateTask", mock.AnythingOfType("*context.timerCtx"), task).Return(expectedError)

		err := usecase.UpdateTask(suite.ctx, task)

		suite.Error(err)
		suite.Equal(expectedError, err)
		mockRepo.AssertExpectations(suite.T())
	})
}

// TestDeleteTaskSuite tests the DeleteTask method
func (suite *TaskUsecaseTestSuite) TestDeleteTaskSuite() {
	suite.Run("Success", func() {
		suite.mockRepo.On("DeleteTask", mock.AnythingOfType("*context.timerCtx"), "task123").Return(nil)

		err := suite.usecase.DeleteTask(suite.ctx, "task123")

		suite.NoError(err)
	})

	suite.Run("EmptyID", func() {
		err := suite.usecase.DeleteTask(suite.ctx, "")

		suite.Error(err)
		suite.Equal("task ID is required", err.Error())
	})

	suite.Run("Error", func() {
		expectedError := errors.New("task not found")
		suite.mockRepo.On("DeleteTask", mock.AnythingOfType("*context.timerCtx"), "nonexistent").Return(expectedError)

		err := suite.usecase.DeleteTask(suite.ctx, "nonexistent")

		suite.Error(err)
		suite.Equal(expectedError, err)
	})
}

// TestContextTimeoutSuite tests context timeout scenarios
func (suite *TaskUsecaseTestSuite) TestContextTimeoutSuite() {
	suite.Run("Timeout", func() {
		// Create usecase with very short timeout
		usecase := NewTaskUsecase(suite.mockRepo, 1*time.Millisecond)

		task := &domain.Task{
			ID:          "task123",
			Title:       "Test Task",
			Description: "Test Description",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
		}

		suite.mockRepo.On("AddTask", mock.AnythingOfType("*context.timerCtx"), task).Return(errors.New("context deadline exceeded"))

		err := usecase.Create(suite.ctx, task)

		suite.Error(err)
		suite.Contains(err.Error(), "context deadline exceeded")
	})
}

// TestTaskUsecaseSuite runs the test suite
func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseTestSuite))
} 
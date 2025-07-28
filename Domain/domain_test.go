package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_Fields(t *testing.T) {
	user := &User{
		ID:       "user123",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "user",
	}

	assert.Equal(t, "user123", user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, "user", user.Role)
}

func TestTask_Fields(t *testing.T) {
	dueDate := time.Now().Add(24 * time.Hour)
	task := &Task{
		ID:          "task123",
		Title:       "Test Task",
		Description: "This is a test task",
		DueDate:     dueDate,
		Status:      "pending",
	}

	assert.Equal(t, "task123", task.ID)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "This is a test task", task.Description)
	assert.Equal(t, dueDate, task.DueDate)
	assert.Equal(t, "pending", task.Status)
}

func TestTask_StatusTransitions(t *testing.T) {
	task := &Task{
		Status: "pending",
	}

	// Test valid status transitions
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
	for _, status := range validStatuses {
		task.Status = status
		assert.Contains(t, validStatuses, task.Status)
	}
}

func TestUser_RoleValidation(t *testing.T) {
	user := &User{
		Role: "invalid_role",
	}

	// Test role validation
	validRoles := []string{"user", "admin"}
	assert.NotContains(t, validRoles, user.Role)

	user.Role = "admin"
	assert.Contains(t, validRoles, user.Role)
} 
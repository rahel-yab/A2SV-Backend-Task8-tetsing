package infrastructure

import (
	"strings"
	"testing"

	"task_manager/domain"

	"github.com/stretchr/testify/suite"
)

// JWTServiceTestSuite is a test suite for JWT service
type JWTServiceTestSuite struct {
	suite.Suite
	jwtService *jwtService
}

// SetupSuite runs once before all tests in the suite
func (suite *JWTServiceTestSuite) SetupSuite() {
	// Suite-level setup if needed
}

// SetupTest runs before each test
func (suite *JWTServiceTestSuite) SetupTest() {
	suite.jwtService = &jwtService{
		secretKey: []byte("test_secret_key"),
	}
}

// TestJWTServiceSuite tests the JWT service functionality
func (suite *JWTServiceTestSuite) TestJWTServiceSuite() {
	suite.Run("GenerateToken_Success", func() {
		// Arrange
		user := &domain.User{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
		}

		// Act
		token, err := suite.jwtService.GenerateToken(user)

		// Assert
		suite.NoError(err)
		suite.NotEmpty(token)
		suite.Contains(token, ".")
		suite.Equal(3, len(strings.Split(token, "."))) // JWT has 3 parts separated by dots
	})

	suite.Run("GenerateToken_EmptySecret", func() {
		// Arrange
		jwtService := NewJWTService("")
		user := &domain.User{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
		}

		// Act
		token, err := jwtService.GenerateToken(user)

		// Assert
		suite.NoError(err)
		suite.NotEmpty(token)
	})

	suite.Run("GenerateToken_NilUser", func() {
		// Act
		token, err := suite.jwtService.GenerateToken(nil)

		// Assert
		suite.Error(err)
		suite.Equal("user cannot be nil", err.Error())
		suite.Empty(token)
	})

	suite.Run("GenerateToken_UserWithEmptyFields", func() {
		// Arrange
		user := &domain.User{
			ID:       "",
			Username: "",
			Email:    "",
			Password: "",
			Role:     "",
		}

		// Act
		token, err := suite.jwtService.GenerateToken(user)

		// Assert
		suite.NoError(err)
		suite.NotEmpty(token)
	})

	suite.Run("GenerateToken_AdminUser", func() {
		// Arrange
		user := &domain.User{
			ID:       "admin123",
			Username: "adminuser",
			Email:    "admin@example.com",
			Password: "hashedpassword",
			Role:     "admin",
		}

		// Act
		token, err := suite.jwtService.GenerateToken(user)

		// Assert
		suite.NoError(err)
		suite.NotEmpty(token)
	})
}

// TestJWTServiceIntegrationSuite tests integration scenarios
func (suite *JWTServiceTestSuite) TestJWTServiceIntegrationSuite() {
	suite.Run("DifferentSecrets", func() {
		// Arrange
		jwtService1 := NewJWTService("secret1")
		jwtService2 := NewJWTService("secret2")
		user := &domain.User{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		}

		// Act
		token1, err1 := jwtService1.GenerateToken(user)
		token2, err2 := jwtService2.GenerateToken(user)

		// Assert
		suite.NoError(err1)
		suite.NoError(err2)
		suite.NotEmpty(token1)
		suite.NotEmpty(token2)
		suite.NotEqual(token1, token2) // Different secrets should produce different tokens
	})

	suite.Run("Integration_WithVariousUsers", func() {
		users := []*domain.User{
			{
				ID:       "user1",
				Username: "user1",
				Email:    "user1@example.com",
				Role:     "user",
			},
			{
				ID:       "user2",
				Username: "user2",
				Email:    "user2@example.com",
				Role:     "admin",
			},
			{
				ID:       "user3",
				Username: "user3",
				Email:    "user3@example.com",
				Role:     "user",
			},
		}

		for i, user := range users {
			suite.Run("User_"+string(rune('1'+i)), func() {
				// Act
				token, err := suite.jwtService.GenerateToken(user)

				// Assert
				suite.NoError(err)
				suite.NotEmpty(token)
				suite.Contains(token, ".")
				suite.Equal(3, len(strings.Split(token, ".")))
			})
		}
	})
}

// TestJWTServiceSuite runs the test suite
func TestJWTServiceSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))
} 
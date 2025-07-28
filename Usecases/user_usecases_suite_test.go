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

// Mock types for the test suite
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) AddUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) IsUsersCollectionEmpty(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) PromoteUserToAdmin(ctx context.Context, identifier string) error {
	args := m.Called(ctx, identifier)
	return args.Error(0)
}

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) CheckPasswordHash(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

// UserUsecaseTestSuite is a test suite for UserUsecase
type UserUsecaseTestSuite struct {
	suite.Suite
	mockUserRepo      *MockUserRepository
	mockPasswordService *MockPasswordService
	mockJWTService    *MockJWTService
	usecase           *UserUsecase
	ctx               context.Context
}

// SetupSuite runs once before all tests in the suite
func (suite *UserUsecaseTestSuite) SetupSuite() {
	suite.ctx = context.Background()
}

// SetupTest runs before each test
func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.mockUserRepo = new(MockUserRepository)
	suite.mockPasswordService = new(MockPasswordService)
	suite.mockJWTService = new(MockJWTService)
	suite.usecase = NewUserUsecase(suite.mockUserRepo, suite.mockPasswordService, suite.mockJWTService, 5*time.Second)
}

// TearDownTest runs after each test
func (suite *UserUsecaseTestSuite) TearDownTest() {
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockPasswordService.AssertExpectations(suite.T())
	suite.mockJWTService.AssertExpectations(suite.T())
}

// TestRegisterUserSuite tests the RegisterUser method
func (suite *UserUsecaseTestSuite) TestRegisterUserSuite() {
	suite.Run("Success_FirstUser", func() {
		username := "testuser"
		email := "test@example.com"
		password := "password123"
		hashedPassword := "hashed_password"

		// Mock repository calls
		suite.mockUserRepo.On("UserExistsByEmail", mock.AnythingOfType("*context.timerCtx"), email).Return(false, nil)
		suite.mockUserRepo.On("UserExistsByUsername", mock.AnythingOfType("*context.timerCtx"), username).Return(false, nil)
		suite.mockUserRepo.On("IsUsersCollectionEmpty", mock.AnythingOfType("*context.timerCtx")).Return(true, nil)
		suite.mockPasswordService.On("HashPassword", password).Return(hashedPassword, nil)
		suite.mockUserRepo.On("AddUser", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.User")).Return(nil)

		role, err := suite.usecase.RegisterUser(suite.ctx, username, email, password)

		suite.NoError(err)
		suite.Equal("admin", role)
	})

	suite.Run("Success_RegularUser", func() {
		// Create new mocks for this specific test to avoid interference
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		usecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService, 5*time.Second)
		
		username := "testuser"
		email := "test@example.com"
		password := "password123"
		hashedPassword := "hashed_password"

		// Mock repository calls
		mockUserRepo.On("UserExistsByEmail", mock.AnythingOfType("*context.timerCtx"), email).Return(false, nil)
		mockUserRepo.On("UserExistsByUsername", mock.AnythingOfType("*context.timerCtx"), username).Return(false, nil)
		mockUserRepo.On("IsUsersCollectionEmpty", mock.AnythingOfType("*context.timerCtx")).Return(false, nil)
		mockPasswordService.On("HashPassword", password).Return(hashedPassword, nil)
		mockUserRepo.On("AddUser", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.User")).Return(nil)

		role, err := usecase.RegisterUser(suite.ctx, username, email, password)

		suite.NoError(err)
		suite.Equal("user", role)
		
		mockUserRepo.AssertExpectations(suite.T())
		mockPasswordService.AssertExpectations(suite.T())
		mockJWTService.AssertExpectations(suite.T())
	})

	suite.Run("EmptyUsername", func() {
		username := ""
		email := "test@example.com"
		password := "password123"

		role, err := suite.usecase.RegisterUser(suite.ctx, username, email, password)

		suite.Error(err)
		suite.Equal("username is required", err.Error())
		suite.Empty(role)
	})

	suite.Run("EmptyEmail", func() {
		username := "testuser"
		email := ""
		password := "password123"

		role, err := suite.usecase.RegisterUser(suite.ctx, username, email, password)

		suite.Error(err)
		suite.Equal("email is required", err.Error())
		suite.Empty(role)
	})

	suite.Run("EmptyPassword", func() {
		username := "testuser"
		email := "test@example.com"
		password := ""

		role, err := suite.usecase.RegisterUser(suite.ctx, username, email, password)

		suite.Error(err)
		suite.Equal("password is required", err.Error())
		suite.Empty(role)
	})

	suite.Run("ShortPassword", func() {
		username := "testuser"
		email := "test@example.com"
		password := "12345" // Less than 6 characters

		role, err := suite.usecase.RegisterUser(suite.ctx, username, email, password)

		suite.Error(err)
		suite.Equal("password must be at least 6 characters long", err.Error())
		suite.Empty(role)
	})

	suite.Run("InvalidEmail", func() {
		username := "testuser"
		email := "invalid-email" // No @ symbol
		password := "password123"

		role, err := suite.usecase.RegisterUser(suite.ctx, username, email, password)

		suite.Error(err)
		suite.Equal("invalid email format", err.Error())
		suite.Empty(role)
	})

	suite.Run("EmailAlreadyExists", func() {
		username := "testuser"
		email := "existing@example.com"
		password := "password123"

		suite.mockUserRepo.On("UserExistsByEmail", mock.AnythingOfType("*context.timerCtx"), email).Return(true, nil)

		role, err := suite.usecase.RegisterUser(suite.ctx, username, email, password)

		suite.Error(err)
		suite.Equal("email already registered", err.Error())
		suite.Empty(role)
	})

	suite.Run("UsernameAlreadyExists", func() {
		username := "existinguser"
		email := "test@example.com"
		password := "password123"

		suite.mockUserRepo.On("UserExistsByEmail", mock.AnythingOfType("*context.timerCtx"), email).Return(false, nil)
		suite.mockUserRepo.On("UserExistsByUsername", mock.AnythingOfType("*context.timerCtx"), username).Return(true, nil)

		role, err := suite.usecase.RegisterUser(suite.ctx, username, email, password)

		suite.Error(err)
		suite.Equal("username already taken", err.Error())
		suite.Empty(role)
	})

	suite.Run("PasswordHashingError", func() {
		// Create new mocks for this specific test to avoid interference
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		usecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService, 5*time.Second)
		
		username := "testuser"
		email := "test@example.com"
		password := "password123"

		mockUserRepo.On("UserExistsByEmail", mock.AnythingOfType("*context.timerCtx"), email).Return(false, nil)
		mockUserRepo.On("UserExistsByUsername", mock.AnythingOfType("*context.timerCtx"), username).Return(false, nil)
		mockUserRepo.On("IsUsersCollectionEmpty", mock.AnythingOfType("*context.timerCtx")).Return(false, nil)
		mockPasswordService.On("HashPassword", password).Return("", errors.New("hashing error"))

		role, err := usecase.RegisterUser(suite.ctx, username, email, password)

		suite.Error(err)
		suite.Equal("hashing error", err.Error())
		suite.Empty(role)
		
		mockUserRepo.AssertExpectations(suite.T())
		mockPasswordService.AssertExpectations(suite.T())
		mockJWTService.AssertExpectations(suite.T())
	})
}

// TestLoginUserSuite tests the LoginUser method
func (suite *UserUsecaseTestSuite) TestLoginUserSuite() {
	suite.Run("Success_WithEmail", func() {
		usernameOrEmail := "test@example.com"
		password := "password123"
		token := "jwt_token"

		user := &domain.User{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed_password",
			Role:     "user",
		}

		suite.mockUserRepo.On("GetUserByEmail", mock.AnythingOfType("*context.timerCtx"), usernameOrEmail).Return(user, nil)
		suite.mockPasswordService.On("CheckPasswordHash", password, user.Password).Return(true)
		suite.mockJWTService.On("GenerateToken", user).Return(token, nil)

		resultToken, role, err := suite.usecase.LoginUser(suite.ctx, usernameOrEmail, password)

		suite.NoError(err)
		suite.Equal(token, resultToken)
		suite.Equal("user", role)
	})

	suite.Run("Success_WithUsername", func() {
		usernameOrEmail := "testuser"
		password := "password123"
		token := "jwt_token"

		user := &domain.User{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed_password",
			Role:     "admin",
		}

		suite.mockUserRepo.On("GetUserByEmail", mock.AnythingOfType("*context.timerCtx"), usernameOrEmail).Return(nil, errors.New("user not found"))
		suite.mockUserRepo.On("GetUserByUsername", mock.AnythingOfType("*context.timerCtx"), usernameOrEmail).Return(user, nil)
		suite.mockPasswordService.On("CheckPasswordHash", password, user.Password).Return(true)
		suite.mockJWTService.On("GenerateToken", user).Return(token, nil)

		resultToken, role, err := suite.usecase.LoginUser(suite.ctx, usernameOrEmail, password)

		suite.NoError(err)
		suite.Equal(token, resultToken)
		suite.Equal("admin", role)
	})

	suite.Run("EmptyUsernameOrEmail", func() {
		usernameOrEmail := ""
		password := "password123"

		resultToken, role, err := suite.usecase.LoginUser(suite.ctx, usernameOrEmail, password)

		suite.Error(err)
		suite.Equal("username or email is required", err.Error())
		suite.Empty(resultToken)
		suite.Empty(role)
	})

	suite.Run("EmptyPassword", func() {
		usernameOrEmail := "testuser"
		password := ""

		resultToken, role, err := suite.usecase.LoginUser(suite.ctx, usernameOrEmail, password)

		suite.Error(err)
		suite.Equal("password is required", err.Error())
		suite.Empty(resultToken)
		suite.Empty(role)
	})

	suite.Run("UserNotFound", func() {
		usernameOrEmail := "nonexistent"
		password := "password123"

		suite.mockUserRepo.On("GetUserByEmail", mock.AnythingOfType("*context.timerCtx"), usernameOrEmail).Return(nil, errors.New("user not found"))
		suite.mockUserRepo.On("GetUserByUsername", mock.AnythingOfType("*context.timerCtx"), usernameOrEmail).Return(nil, errors.New("user not found"))

		resultToken, role, err := suite.usecase.LoginUser(suite.ctx, usernameOrEmail, password)

		suite.Error(err)
		suite.Equal("invalid email/username or password", err.Error())
		suite.Empty(resultToken)
		suite.Empty(role)
	})

	suite.Run("WrongPassword", func() {
		usernameOrEmail := "test@example.com"
		password := "wrongpassword"

		user := &domain.User{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed_password",
			Role:     "user",
		}

		suite.mockUserRepo.On("GetUserByEmail", mock.AnythingOfType("*context.timerCtx"), usernameOrEmail).Return(user, nil)
		suite.mockPasswordService.On("CheckPasswordHash", password, user.Password).Return(false)

		resultToken, role, err := suite.usecase.LoginUser(suite.ctx, usernameOrEmail, password)

		suite.Error(err)
		suite.Equal("invalid email/username or password", err.Error())
		suite.Empty(resultToken)
		suite.Empty(role)
	})

	suite.Run("JWTGenerationError", func() {
		// Create new mocks for this specific test
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		usecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService, 5*time.Second)
		
		usernameOrEmail := "test@example.com"
		password := "password123"

		user := &domain.User{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed_password",
			Role:     "user",
		}

		mockUserRepo.On("GetUserByEmail", mock.AnythingOfType("*context.timerCtx"), usernameOrEmail).Return(user, nil)
		mockPasswordService.On("CheckPasswordHash", password, user.Password).Return(true)
		mockJWTService.On("GenerateToken", user).Return("", errors.New("JWT generation error"))

		resultToken, role, err := usecase.LoginUser(suite.ctx, usernameOrEmail, password)

		suite.Error(err)
		suite.Equal("JWT generation error", err.Error())
		suite.Empty(resultToken)
		suite.Empty(role)
		
		mockUserRepo.AssertExpectations(suite.T())
		mockPasswordService.AssertExpectations(suite.T())
		mockJWTService.AssertExpectations(suite.T())
	})
}

// TestPromoteUserToAdminSuite tests the PromoteUserToAdmin method
func (suite *UserUsecaseTestSuite) TestPromoteUserToAdminSuite() {
	suite.Run("Success", func() {
		identifier := "testuser"

		suite.mockUserRepo.On("PromoteUserToAdmin", mock.AnythingOfType("*context.timerCtx"), identifier).Return(nil)

		err := suite.usecase.PromoteUserToAdmin(suite.ctx, identifier)

		suite.NoError(err)
	})

	suite.Run("EmptyIdentifier", func() {
		identifier := ""

		err := suite.usecase.PromoteUserToAdmin(suite.ctx, identifier)

		suite.Error(err)
		suite.Equal("identifier is required", err.Error())
	})

	suite.Run("Error", func() {
		identifier := "nonexistent"

		suite.mockUserRepo.On("PromoteUserToAdmin", mock.AnythingOfType("*context.timerCtx"), identifier).Return(errors.New("user not found"))

		err := suite.usecase.PromoteUserToAdmin(suite.ctx, identifier)

		suite.Error(err)
		suite.Equal("user not found", err.Error())
	})
}

// TestUserUsecaseSuite runs the test suite
func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
} 
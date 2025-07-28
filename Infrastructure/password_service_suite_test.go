package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// PasswordServiceTestSuite is a test suite for password service
type PasswordServiceTestSuite struct {
	suite.Suite
	passwordService *passwordService
}

// SetupSuite runs once before all tests in the suite
func (suite *PasswordServiceTestSuite) SetupSuite() {
	// Suite-level setup if needed
}

// SetupTest runs before each test
func (suite *PasswordServiceTestSuite) SetupTest() {
	suite.passwordService = &passwordService{}
}

// TestPasswordServiceSuite tests the password service functionality
func (suite *PasswordServiceTestSuite) TestPasswordServiceSuite() {
	suite.Run("HashPassword_Success", func() {
		// Arrange
		password := "testpassword123"

		// Act
		hash, err := suite.passwordService.HashPassword(password)

		// Assert
		suite.NoError(err)
		suite.NotEmpty(hash)
		suite.NotEqual(password, hash)
		suite.Len(hash, 60) // bcrypt hash length
	})

	suite.Run("HashPassword_EmptyPassword", func() {
		// Arrange
		password := ""

		// Act
		hash, err := suite.passwordService.HashPassword(password)

		// Assert
		suite.NoError(err)
		suite.NotEmpty(hash)
		suite.NotEqual(password, hash)
	})

	suite.Run("CheckPasswordHash_Success", func() {
		// Arrange
		password := "testpassword123"
		hash, _ := suite.passwordService.HashPassword(password)

		// Act
		isValid := suite.passwordService.CheckPasswordHash(password, hash)

		// Assert
		suite.True(isValid)
	})

	suite.Run("CheckPasswordHash_WrongPassword", func() {
		// Arrange
		correctPassword := "testpassword123"
		wrongPassword := "wrongpassword"
		hash, _ := suite.passwordService.HashPassword(correctPassword)

		// Act
		isValid := suite.passwordService.CheckPasswordHash(wrongPassword, hash)

		// Assert
		suite.False(isValid)
	})

	suite.Run("CheckPasswordHash_EmptyPassword", func() {
		// Arrange
		password := ""
		hash, _ := suite.passwordService.HashPassword(password)

		// Act
		isValid := suite.passwordService.CheckPasswordHash(password, hash)

		// Assert
		suite.True(isValid)
	})

	suite.Run("CheckPasswordHash_EmptyHash", func() {
		// Arrange
		password := "testpassword123"
		hash := ""

		// Act
		isValid := suite.passwordService.CheckPasswordHash(password, hash)

		// Assert
		suite.False(isValid)
	})

	suite.Run("CheckPasswordHash_InvalidHash", func() {
		// Arrange
		password := "testpassword123"
		invalidHash := "invalid_hash_format"

		// Act
		isValid := suite.passwordService.CheckPasswordHash(password, invalidHash)

		// Assert
		suite.False(isValid)
	})
}

// TestPasswordServiceIntegrationSuite tests integration scenarios
func (suite *PasswordServiceTestSuite) TestPasswordServiceIntegrationSuite() {
	suite.Run("HashAndVerify_Integration", func() {
		testPasswords := []string{
			"simple",
			"complex_password_123!@#",
			"very_long_password_with_special_characters_!@#$%^&*()_+-=[]{}|;:,.<>?",
			"",
			"1234567890",
		}

		for _, password := range testPasswords {
			suite.Run("Password_"+password, func() {
				// Act - Hash the password
				hash, err := suite.passwordService.HashPassword(password)

				// Assert - Hash should be successful
				suite.NoError(err)
				suite.NotEmpty(hash)
				suite.NotEqual(password, hash)

				// Act - Verify the password
				isValid := suite.passwordService.CheckPasswordHash(password, hash)

				// Assert - Verification should be successful
				suite.True(isValid)

				// Act - Verify wrong password
				wrongPassword := password + "_wrong"
				isValidWrong := suite.passwordService.CheckPasswordHash(wrongPassword, hash)

				// Assert - Wrong password should fail
				suite.False(isValidWrong)
			})
		}
	})
}

// TestPasswordServiceSuite runs the test suite
func TestPasswordServiceSuite(t *testing.T) {
	suite.Run(t, new(PasswordServiceTestSuite))
} 
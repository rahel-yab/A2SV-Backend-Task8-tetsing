package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"task_manager/domain"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// AuthMiddlewareTestSuite is a test suite for authentication middleware
type AuthMiddlewareTestSuite struct {
	suite.Suite
	jwtSecret []byte
	router    *gin.Engine
}

// SetupSuite runs once before all tests in the suite
func (suite *AuthMiddlewareTestSuite) SetupSuite() {
	suite.jwtSecret = []byte("test_secret")
}

// SetupTest runs before each test
func (suite *AuthMiddlewareTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
}

// setupTestRouter creates a test router with auth middleware
func (suite *AuthMiddlewareTestSuite) setupTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(AuthMiddleware(suite.jwtSecret))
	return router
}

// TestAuthMiddlewareSuite tests the AuthMiddleware functionality
func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareSuite() {
	suite.Run("ValidToken", func() {
		// Create a valid JWT token
		user := &domain.User{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		}
		
		jwtService := NewJWTService(string(suite.jwtSecret))
		token, err := jwtService.GenerateToken(user)
		suite.NoError(err)

		router := suite.setupTestRouter()
		router.GET("/test", func(c *gin.Context) {
			claims, exists := c.Get("claims")
			suite.True(exists)
			c.JSON(http.StatusOK, gin.H{"message": "success", "claims": claims})
		})

		// Act
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		suite.Equal(http.StatusOK, w.Code)
		suite.Contains(w.Body.String(), "success")
	})

	suite.Run("MissingAuthorizationHeader", func() {
		router := suite.setupTestRouter()
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Act
		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		suite.Equal(http.StatusUnauthorized, w.Code)
		suite.Contains(w.Body.String(), "Authorization header is required")
	})

	suite.Run("InvalidAuthorizationFormat", func() {
		router := suite.setupTestRouter()
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		testCases := []struct {
			header       string
			expectedMsg  string
		}{
			{"InvalidToken", "Invalid authorization header"},
			{"Bearer", "Invalid authorization header"},
			{"Basic token", "Invalid authorization header"},
			{"Bearer ", "Invalid JWT"},
			{"", "Authorization header is required"},
		}

		for _, tc := range testCases {
			suite.Run("Header_"+tc.header, func() {
				// Act
				req, _ := http.NewRequest("GET", "/test", nil)
				if tc.header != "" {
					req.Header.Set("Authorization", tc.header)
				}
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Assert
				suite.Equal(http.StatusUnauthorized, w.Code)
				suite.Contains(w.Body.String(), tc.expectedMsg)
			})
		}
	})

	suite.Run("InvalidToken", func() {
		router := suite.setupTestRouter()
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Act
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		suite.Equal(http.StatusUnauthorized, w.Code)
		suite.Contains(w.Body.String(), "Invalid JWT")
	})

	suite.Run("ExpiredToken", func() {
		// Create an expired JWT token
		claims := jwt.MapClaims{
			"user_id":  "user123",
			"username": "testuser",
			"email":    "test@example.com",
			"role":     "user",
			"exp":      time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(suite.jwtSecret)
		suite.NoError(err)

		router := suite.setupTestRouter()
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Act
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		suite.Equal(http.StatusUnauthorized, w.Code)
		suite.Contains(w.Body.String(), "Invalid JWT")
	})

	suite.Run("WrongSecret", func() {
		// Create a token with wrong secret
		user := &domain.User{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		}
		
		jwtService := NewJWTService("wrong_secret")
		token, err := jwtService.GenerateToken(user)
		suite.NoError(err)

		router := suite.setupTestRouter()
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Act
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		suite.Equal(http.StatusUnauthorized, w.Code)
		suite.Contains(w.Body.String(), "Invalid JWT")
	})
}

// TestAdminOnlySuite tests the AdminOnly middleware functionality
func (suite *AuthMiddlewareTestSuite) TestAdminOnlySuite() {
	suite.Run("ValidAdminUser", func() {
		// Create a valid admin JWT token
		user := &domain.User{
			ID:       "admin123",
			Username: "adminuser",
			Email:    "admin@example.com",
			Role:     "admin",
		}
		
		jwtService := NewJWTService(string(suite.jwtSecret))
		token, err := jwtService.GenerateToken(user)
		suite.NoError(err)

		router := suite.setupTestRouter()
		router.GET("/admin", AdminOnly(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})

		// Act
		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		suite.Equal(http.StatusOK, w.Code)
		suite.Contains(w.Body.String(), "admin access granted")
	})

	suite.Run("RegularUser", func() {
		// Create a valid user JWT token (not admin)
		user := &domain.User{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		}
		
		jwtService := NewJWTService(string(suite.jwtSecret))
		token, err := jwtService.GenerateToken(user)
		suite.NoError(err)

		router := suite.setupTestRouter()
		router.GET("/admin", AdminOnly(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})

		// Act
		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		suite.Equal(http.StatusForbidden, w.Code)
		suite.Contains(w.Body.String(), "Admin access required")
	})

	suite.Run("NoClaims", func() {
		router := gin.New()
		router.GET("/admin", AdminOnly(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})

		// Act
		req, _ := http.NewRequest("GET", "/admin", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		suite.Equal(http.StatusForbidden, w.Code)
		suite.Contains(w.Body.String(), "Admin access required")
	})
}

// TestAuthMiddlewareIntegrationSuite tests integration scenarios
func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareIntegrationSuite() {
	suite.Run("IntegrationScenarios", func() {
		// Create tokens for different user types
		adminUser := &domain.User{
			ID:       "admin123",
			Username: "adminuser",
			Email:    "admin@example.com",
			Role:     "admin",
		}
		
		regularUser := &domain.User{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		}
		
		jwtService := NewJWTService(string(suite.jwtSecret))
		adminToken, _ := jwtService.GenerateToken(adminUser)
		userToken, _ := jwtService.GenerateToken(regularUser)

		router := suite.setupTestRouter()
		router.GET("/public", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "public access"})
		})
		router.GET("/admin", AdminOnly(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access"})
		})

		testCases := []struct {
			name           string
			endpoint       string
			token          string
			expectedStatus int
			expectedBody   string
		}{
			{"Public endpoint with admin token", "/public", adminToken, http.StatusOK, "public access"},
			{"Public endpoint with user token", "/public", userToken, http.StatusOK, "public access"},
			{"Admin endpoint with admin token", "/admin", adminToken, http.StatusOK, "admin access"},
			{"Admin endpoint with user token", "/admin", userToken, http.StatusForbidden, "Admin access required"},
			{"Admin endpoint without token", "/admin", "", http.StatusUnauthorized, "Authorization header is required"},
		}

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				// Act
				req, _ := http.NewRequest("GET", tc.endpoint, nil)
				if tc.token != "" {
					req.Header.Set("Authorization", "Bearer "+tc.token)
				}
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Assert
				suite.Equal(tc.expectedStatus, w.Code)
				suite.Contains(w.Body.String(), tc.expectedBody)
			})
		}
	})
}

// TestAuthMiddlewareSuite runs the test suite
func TestAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
} 
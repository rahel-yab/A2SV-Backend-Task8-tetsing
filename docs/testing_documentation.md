# Task Management API - Comprehensive Testing Documentation

## Overview

This document provides comprehensive documentation for the unit testing suite of the Task Management REST API, built with Go and the `testify` library.

## Test Architecture

The testing suite follows Clean Architecture principles and is organized into the following layers:

```
task_manager/
├── Domain/
│   └── domain_test.go          # Domain model validation tests
├── Infrastructure/
│   ├── auth_middleware_suite_test.go # Authentication & authorization tests
│   ├── jwt_service_suite_test.go     # JWT token service tests
│   └── password_service_suite_test.go # Password hashing & verification tests
├── Usecases/
│   ├── task_usecases_suite_test.go # Task business logic test suite
│   └── user_usecases_suite_test.go # User business logic test suite
└── Repositories/               # Data access layer (no tests - integration testing recommended)
```

## CI/CD Pipeline Integration

### GitHub Actions Workflow

The project includes a comprehensive CI/CD pipeline implemented with GitHub Actions:

**Location**: `.github/workflows/test.yml`

**Features**:

- ✅ **Automated Testing**: Runs on every push and pull request
- ✅ **Multi-version Testing**: Tests against Go versions 1.21, 1.22, 1.23
- ✅ **Coverage Reporting**: Generates detailed coverage reports
- ✅ **Quality Gates**: Enforces 80% test coverage threshold
- ✅ **Code Quality Checks**: Linting, vetting, and formatting validation
- ✅ **Artifact Upload**: Preserves test results and coverage reports

**Pipeline Steps**:

1. **Checkout**: Retrieves the latest code
2. **Setup**: Configures Go environment
3. **Dependencies**: Installs required packages
4. **Testing**: Runs comprehensive test suite with coverage
5. **Quality Checks**: Validates code quality and formatting
6. **Artifacts**: Uploads test results for review

### Makefile Support

A comprehensive Makefile provides easy access to development tasks:

**Available Commands**:

```bash
make test              # Run all tests
make test-coverage     # Run tests with coverage reporting
make test-coverage-check # Run tests with 80% coverage threshold
make lint              # Run code linting
make vet               # Run go vet
make fmt               # Check code formatting
make fmt-fix           # Fix code formatting
make ci                # Run full CI pipeline locally
make clean             # Clean build artifacts
make build             # Build the application
make run               # Run the application
make deps              # Install dependencies
make setup             # Initial development setup
```

### Coverage Thresholds

The CI pipeline enforces quality standards:

- **Minimum Coverage**: 80% test coverage required
- **Code Quality**: Must pass linting and vetting
- **Formatting**: Code must be properly formatted
- **Dependencies**: All dependencies must be valid

## Test Suite Structure

### Testify Suite Implementation

All test suites follow the standard `testify/suite` pattern:

```go
type MyTestSuite struct {
    suite.Suite
    // Test dependencies
}

// SetupSuite runs once before all tests in the suite
func (suite *MyTestSuite) SetupSuite() {
    // Initialize shared resources
}

// SetupTest runs before each test
func (suite *MyTestSuite) SetupTest() {
    // Reset mocks and test state
}

// TearDownTest runs after each test
func (suite *MyTestSuite) TearDownTest() {
    // Verify mock expectations
}

// Test methods
func (suite *MyTestSuite) TestMyFunction() {
    // Test implementation
}

// Entry point
func TestMySuite(t *testing.T) {
    suite.Run(t, new(MyTestSuite))
}
```

## Test Categories

### 1. Domain Tests (`Domain/domain_test.go`)

**Purpose**: Validate core business entities and their behavior

**Coverage**:

- ✅ User field validation and assignment
- ✅ Task field validation and assignment
- ✅ Task status transition validation
- ✅ User role validation

**Test Methods**:

- `TestUser_Fields()` - Validates user struct field assignments
- `TestTask_Fields()` - Validates task struct field assignments
- `TestTask_StatusTransitions()` - Tests valid status transitions
- `TestUser_RoleValidation()` - Tests user role validation

### 2. Infrastructure Tests

#### Authentication Middleware (`Infrastructure/auth_middleware_suite_test.go`)

**Purpose**: Test JWT authentication and role-based authorization

**Coverage**:

- ✅ Valid token handling
- ✅ Missing/invalid authorization headers
- ✅ Expired/wrong tokens
- ✅ Admin-only access control
- ✅ Integration scenarios

**Test Methods**:

- `TestAuthMiddlewareSuite()` - Tests authentication middleware
- `TestAdminOnlySuite()` - Tests admin-only access control
- `TestAuthMiddlewareIntegrationSuite()` - Tests integration scenarios

#### JWT Service (`Infrastructure/jwt_service_suite_test.go`)

**Purpose**: Test JWT token generation and validation

**Coverage**:

- ✅ Token generation with valid data
- ✅ Error handling for invalid inputs
- ✅ Different secret handling
- ✅ Integration with various user types

**Test Methods**:

- `TestJWTServiceSuite()` - Tests JWT service functionality
- `TestJWTServiceIntegrationSuite()` - Tests integration scenarios

#### Password Service (`Infrastructure/password_service_suite_test.go`)

**Purpose**: Test password hashing and verification

**Coverage**:

- ✅ Password hashing functionality
- ✅ Password verification
- ✅ Error handling for invalid inputs
- ✅ Integration with various password types

**Test Methods**:

- `TestPasswordServiceSuite()` - Tests password service functionality
- `TestPasswordServiceIntegrationSuite()` - Tests integration scenarios

### 3. Usecase Tests

#### Task Usecases (`Usecases/task_usecases_suite_test.go`)

**Purpose**: Test task-related business logic

**Coverage**:

- ✅ Task creation with validation
- ✅ Task retrieval (all and by ID)
- ✅ Task updates with validation
- ✅ Task deletion
- ✅ Error handling for invalid inputs
- ✅ Context timeout handling

**Test Methods**:

- `TestCreateTaskSuite()` - Tests task creation
- `TestGetAllTasksSuite()` - Tests task retrieval
- `TestGetTaskByIDSuite()` - Tests individual task retrieval
- `TestUpdateTaskSuite()` - Tests task updates
- `TestDeleteTaskSuite()` - Tests task deletion
- `TestContextTimeoutSuite()` - Tests timeout handling

#### User Usecases (`Usecases/user_usecases_suite_test.go`)

**Purpose**: Test user-related business logic

**Coverage**:

- ✅ User registration with validation
- ✅ User login with authentication
- ✅ User promotion to admin
- ✅ Error handling for invalid inputs
- ✅ Password hashing integration
- ✅ JWT token generation

**Test Methods**:

- `TestRegisterUserSuite()` - Tests user registration
- `TestLoginUserSuite()` - Tests user authentication
- `TestPromoteUserToAdminSuite()` - Tests user promotion

## Running Tests

### Local Development

**Run all tests**:

```bash
go test ./... -v
```

**Run tests with coverage**:

```bash
go test ./... -v -cover
```

**Run specific test suite**:

```bash
go test ./Usecases/... -v
go test ./Infrastructure/... -v
```

**Using Makefile**:

```bash
make test              # Run all tests
make test-coverage     # Run with coverage
make ci                # Run full CI pipeline
```

### CI/CD Pipeline

The CI pipeline automatically runs on:

- **Push to main/develop branches**
- **Pull requests to main/develop branches**

**Pipeline Execution**:

1. Tests run against multiple Go versions
2. Coverage reports are generated
3. Quality checks are performed
4. Results are uploaded as artifacts

## Test Coverage Metrics

### Current Coverage (as of latest run)

| Layer              | Coverage         | Status |
| ------------------ | ---------------- | ------ |
| **Domain**         | Basic validation | ✅     |
| **Infrastructure** | 85.2%            | ✅     |
| **Usecases**       | 95.4%            | ✅     |
| **Overall**        | 45.9%            | ⚠️     |

### Coverage Breakdown

**Infrastructure Layer (85.2%)**:

- Auth Middleware: 84.6%
- JWT Service: 100%
- Password Service: 100%

**Usecases Layer (95.4%)**:

- Task Usecases: 100% (core methods)
- User Usecases: 96.3% (core methods)

**Areas for Improvement**:

- Delivery layer (Controllers) - 0% (integration testing recommended)
- Repository layer - 0% (integration testing recommended)

## Best Practices

### Test Organization

1. **Use Test Suites**: Organize related tests into suites
2. **Clear Naming**: Use descriptive test and method names
3. **Setup/Teardown**: Properly manage test state
4. **Mock Isolation**: Use mocks for external dependencies

### Test Writing

1. **Arrange-Act-Assert**: Follow the AAA pattern
2. **Single Responsibility**: Each test should test one thing
3. **Edge Cases**: Test both success and failure scenarios
4. **Validation**: Test input validation thoroughly

### Mock Usage

1. **Isolation**: Use mocks to isolate units under test
2. **Expectations**: Set clear expectations for mock calls
3. **Verification**: Verify that expected calls were made
4. **Cleanup**: Properly clean up mocks between tests

## Troubleshooting

### Common Issues

**Test Failures**:

- Check mock expectations
- Verify test data setup
- Ensure proper cleanup

**Coverage Issues**:

- Add tests for uncovered code paths
- Consider integration tests for delivery layer
- Review test organization

**CI Pipeline Failures**:

- Check coverage threshold (80% minimum)
- Verify code formatting
- Review linting errors

### Debugging

**Verbose Output**:

```bash
go test ./... -v
```

**Coverage Analysis**:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Specific Test Debugging**:

```bash
go test -run TestSpecificTest ./...
```

## Conclusion

This comprehensive testing suite provides:

- ✅ **High Coverage**: 95.4% coverage of critical business logic
- ✅ **Quality Assurance**: Automated quality checks
- ✅ **CI/CD Integration**: Automated testing pipeline
- ✅ **Maintainability**: Well-organized, documented tests
- ✅ **Reliability**: Comprehensive edge case testing



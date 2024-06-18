package middleware

import (
	"errors"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/elina-chertova/auth-keeper.git/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockUserRepo is a mock implementation of the UserRepo interface
type MockUserRepo struct {
	mockGetUserByUsername func(username string) (*models.User, error)
	mockCreateUser        func(user *models.User) error
}

func (m *MockUserRepo) GetUserByUsername(username string) (*models.User, error) {
	return m.mockGetUserByUsername(username)
}

func (m *MockUserRepo) CreateUser(user *models.User) error {
	return m.mockCreateUser(user)
}

func TestLoadPersonalKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPersonalKey := []byte("mockPersonalKey")
	encryptedMockPersonalKey, _ := security.EncryptPersonalKey(mockPersonalKey)

	tests := []struct {
		name           string
		mockUserRepo   *MockUserRepo
		contextSetup   func(ctx *gin.Context)
		expectedStatus int
		expectedKey    interface{}
	}{
		{
			name: "Valid user and personal key",
			mockUserRepo: &MockUserRepo{
				mockGetUserByUsername: func(username string) (*models.User, error) {
					return &models.User{PersonalKey: encryptedMockPersonalKey}, nil
				},
			},
			contextSetup: func(ctx *gin.Context) {
				ctx.Set("userID", "validUser")
			},
			expectedStatus: http.StatusOK,
			expectedKey:    mockPersonalKey,
		},
		{
			name: "User ID not found in context",
			mockUserRepo: &MockUserRepo{
				mockGetUserByUsername: func(username string) (*models.User, error) {
					return nil, nil
				},
			},
			contextSetup: func(ctx *gin.Context) {
				// Do not set "userID"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedKey:    nil,
		},
		{
			name: "Failed to get user",
			mockUserRepo: &MockUserRepo{
				mockGetUserByUsername: func(username string) (*models.User, error) {
					return nil, errors.New("user not found")
				},
			},
			contextSetup: func(ctx *gin.Context) {
				ctx.Set("userID", "invalidUser")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedKey:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				recorder := httptest.NewRecorder()
				ctx, r := gin.CreateTestContext(recorder)

				tt.contextSetup(ctx)

				r.Use(
					func(c *gin.Context) {
						for k, v := range ctx.Keys {
							c.Set(k, v)
						}
						LoadPersonalKey(tt.mockUserRepo)(c)
					},
				)
				r.GET(
					"/", func(c *gin.Context) {
						c.Status(http.StatusOK)
					},
				)

				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				r.ServeHTTP(recorder, req)

				assert.Equal(t, tt.expectedStatus, recorder.Code)

			},
		)
	}
}

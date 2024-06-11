package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/elina-chertova/auth-keeper.git/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

type MockSecurity struct {
	mock.Mock
}

func (m *MockSecurity) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockSecurity) CheckPasswordHash(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

func (m *MockSecurity) GenerateToken(username string) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}

func setupRouter(userRepo repository.UserRepo, security *MockSecurity) *gin.Engine {
	handler := NewUserHandler(
		userRepo,
	)
	router := gin.New()
	router.POST("/register", handler.Register())
	router.POST("/signup", handler.Signup())
	return router
}

func TestUserHandler_Register(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockSecurity := new(MockSecurity)
	router := setupRouter(mockRepo, mockSecurity)

	t.Run(
		"Successful Registration", func(t *testing.T) {
			user := &models.User{
				Username: "test_user",
				Password: "password",
				Email:    "test@example.com",
			}

			mockRepo.On("CreateUser", mock.Anything).Return(nil).Once()
			mockSecurity.On("HashPassword", "password").Return("hashed_password", nil).Once()
			mockSecurity.On("GenerateToken", "test_user").Return("test_token", nil).Once()

			body, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "User has been created", response["message"])

			userResponse := response["user"].(map[string]interface{})
			assert.Equal(t, "test_user", userResponse["username"])
			assert.Equal(t, "***", userResponse["password"])
			assert.Equal(t, "test@example.com", userResponse["email"])
			assert.Equal(t, "Kioq", userResponse["personal_key"])

		},
	)

}

func TestUserHandler_Signup(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockSecurity := new(MockSecurity)
	router := setupRouter(mockRepo, mockSecurity)

	t.Run(
		"Successful Signup", func(t *testing.T) {
			user := &models.User{
				Username: "test_user",
				Password: "password",
			}

			dbUser := &models.User{
				Username: "test_user",
				Password: "hashed_password",
			}

			mockRepo.On("GetUserByUsername", "test_user").Return(dbUser, nil).Once()
			mockSecurity.On("CheckPasswordHash", "password", "hashed_password").Return(true).Once()
			mockSecurity.On("GenerateToken", "test_user").Return("test_token", nil).Once()

			body, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "User has logged in", response["message"])

		},
	)

	t.Run(
		"Signup with Invalid JSON", func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		},
	)
}

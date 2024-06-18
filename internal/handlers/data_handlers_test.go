package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/elina-chertova/auth-keeper.git/internal/security"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockCreditCardRepo struct {
	SaveNewCreditCardFunc func(card *models.CreditCard) error
	GetCreditCardListFunc func(userID string) ([]*models.CreditCard, error)
}

func (m *MockCreditCardRepo) SaveNewCreditCard(card *models.CreditCard) error {
	if m.SaveNewCreditCardFunc != nil {
		return m.SaveNewCreditCardFunc(card)
	}
	return nil
}

func (m *MockCreditCardRepo) GetCreditCardList(userID string) ([]*models.CreditCard, error) {
	if m.GetCreditCardListFunc != nil {
		return m.GetCreditCardListFunc(userID)
	}
	return nil, nil
}

func TestDataHandler_AddCreditCardHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	personalKey := []byte("1234567890123456")

	validCreditCard := models.CreditCard{
		UserID:     "test_user",
		CardNumber: "1234-5678-9876-5432",
		ExpiryDate: "12/24",
		CVV:        "123",
		CardHolder: "John Doe",
		Metadata:   "metadata",
	}
	validCreditCardJSON, _ := json.Marshal(validCreditCard)
	encryptedData, _ := security.EncryptData(validCreditCardJSON, personalKey)
	validEncryptedData := base64.StdEncoding.EncodeToString(encryptedData)

	tests := []struct {
		name                 string
		setupContext         func(ctx *gin.Context)
		requestBody          string
		mockSaveCreditCard   func(card *models.CreditCard) error
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Valid Request",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("userID", "test_user")
				ctx.Set("personalKey", personalKey)
			},
			requestBody: validEncryptedData,
			mockSaveCreditCard: func(card *models.CreditCard) error {
				return nil
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"message":"Credit Card has been added","status":201}`,
		},
		{
			name: "No userID in Context",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("personalKey", personalKey)
			},
			requestBody:          validEncryptedData,
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"Unauthorized"}`,
		},
		{
			name: "Invalid JSON Body",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("userID", "test_user")
				ctx.Set("personalKey", personalKey)
			},
			requestBody:          "invalid_json",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"json: cannot unmarshal string into Go value of type struct { Data string \"json:\\\"data\\\"\" }"}`,
		},
		{
			name: "Base64 Decode Failure",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("userID", "test_user")
				ctx.Set("personalKey", personalKey)
			},
			requestBody:          `{"data":"invalid_base64"}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Failed to decode base64 data"}`,
		},
		{
			name: "No personalKey in Context",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("userID", "test_user")
			},
			requestBody:          validEncryptedData,
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Personal key not found"}`,
		},
		{
			name: "Data Decryption Failure",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("userID", "test_user")
				ctx.Set("personalKey", []byte("wrong_key"))
			},
			requestBody:          validEncryptedData,
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Failed to decrypt data"}`,
		},
		{
			name: "Unmarshal Decrypted Data Failure",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("userID", "test_user")
				ctx.Set("personalKey", personalKey)
			},
			requestBody:          `{"data":"` + base64.StdEncoding.EncodeToString([]byte("invalid_json")) + `"}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Failed to unmarshal decrypted data"}`,
		},
		{
			name: "Save Credit Card Failure",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("userID", "test_user")
				ctx.Set("personalKey", personalKey)
			},
			requestBody: validEncryptedData,
			mockSaveCreditCard: func(card *models.CreditCard) error {
				return errors.New("failed to save credit card")
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"failed to save credit card"}`,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				mockCCRepo := &MockCreditCardRepo{
					SaveNewCreditCardFunc: tt.mockSaveCreditCard,
				}

				dh := &DataHandler{
					cc: mockCCRepo,
				}

				router := gin.New()
				router.Use(
					func(ctx *gin.Context) {
						tt.setupContext(ctx)
						ctx.Next()
					},
				)
				router.POST("/api/data/add-card", dh.AddCreditCardHandler())

				reqBody := `{"data":"` + tt.requestBody + `"}`

				req, _ := http.NewRequest(
					"POST",
					"/api/data/add-card",
					bytes.NewBufferString(reqBody),
				)
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, tt.expectedStatusCode, w.Code)
			},
		)
	}
}

func TestDataHandler_GetCreditCardHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	personalKey := []byte("1234567890123456")

	validCreditCard := models.CreditCard{
		UserID:     "test_user",
		CardNumber: "1234-5678-9876-5432",
		ExpiryDate: "12/24",
		CVV:        "123",
		CardHolder: "John Doe",
		Metadata:   "metadata",
	}
	validCreditCardList := []*models.CreditCard{&validCreditCard}
	validCreditCardListJSON, _ := json.Marshal(validCreditCardList)
	encryptedData, _ := security.EncryptData(validCreditCardListJSON, personalKey)
	validEncryptedData := base64.StdEncoding.EncodeToString(encryptedData)

	token, _ := security.GenerateToken("test_user")

	tests := []struct {
		name                  string
		setupContext          func(ctx *gin.Context)
		mockGetCreditCardList func(userID string) ([]*models.CreditCard, error)
		expectedStatusCode    int
		expectedResponseBody  string
	}{
		{
			name: "Valid Request",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("userID", "test_user")
				ctx.Set("personalKey", personalKey)
				ctx.Set("token", token)
			},
			mockGetCreditCardList: func(userID string) ([]*models.CreditCard, error) {
				return validCreditCardList, nil
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"Credit card list","body":"` + validEncryptedData + `"}`,
		},

		{
			name: "Failed to get user",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("userID", "test_user")
				ctx.Set("personalKey", personalKey)
				ctx.Set("token", token)
			},
			mockGetCreditCardList: func(userID string) ([]*models.CreditCard, error) {
				return nil, errors.New("failed to get user")
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"failed to get user"}`,
		},
		{
			name: "No personalKey in Context",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("userID", "test_user")
				ctx.Set("token", token)
			},
			mockGetCreditCardList: func(userID string) ([]*models.CreditCard, error) {
				return validCreditCardList, nil
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Personal key not found"}`,
		},
		{
			name: "Failed to encrypt data",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("userID", "test_user")
				ctx.Set("personalKey", []byte("wrong_key"))
				ctx.Set("token", token)
			},
			mockGetCreditCardList: func(userID string) ([]*models.CreditCard, error) {
				return validCreditCardList, nil
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Failed to encrypt data"}`,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				mockCCRepo := &MockCreditCardRepo{
					GetCreditCardListFunc: tt.mockGetCreditCardList,
				}

				dh := &DataHandler{
					cc: mockCCRepo,
				}

				router := gin.New()
				router.Use(
					func(ctx *gin.Context) {
						tt.setupContext(ctx)
						ctx.Next()
					},
				)
				router.GET("/api/data/get-card", dh.GetCreditCardHandler())

				req, _ := http.NewRequest("GET", "/api/data/get-card", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, tt.expectedStatusCode, w.Code)
			},
		)
	}
}

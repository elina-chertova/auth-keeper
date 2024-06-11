package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elina-chertova/auth-keeper.git/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	router := gin.New()
	router.Use(JWTAuth())

	router.GET(
		"/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		},
	)

	t.Run(
		"No Authorization Header", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		},
	)

	t.Run(
		"Invalid Authorization Header Format", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", "InvalidToken")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		},
	)

	t.Run(
		"Valid Authorization Header", func(t *testing.T) {
			token, _ := security.GenerateToken("test_user")
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		},
	)
}

func TestExtractUserID(t *testing.T) {
	router := gin.New()
	router.Use(JWTAuth())
	router.Use(ExtractUserID())

	router.GET(
		"/test", func(c *gin.Context) {
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "userID not found"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"userID": userID})
		},
	)

	t.Run(
		"Valid Token", func(t *testing.T) {
			token, _ := security.GenerateToken("test_user")
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, `{"userID":"test_user"}`, w.Body.String())
		},
	)

	t.Run(
		"Invalid Token", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", "Bearer invalid_token")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		},
	)

	t.Run(
		"No Token", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		},
	)
}

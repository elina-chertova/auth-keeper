package security

import (
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"

	"github.com/elina-chertova/auth-keeper.git/internal/config"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.SecretKey = "secret_key"
}

func TestGenerateToken(t *testing.T) {
	username := "test_user"

	tokenString, err := GenerateToken(username)
	assert.Nil(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestValidateToken(t *testing.T) {
	username := "test_user"

	tokenString, err := GenerateToken(username)
	assert.Nil(t, err)
	err = ValidateToken(tokenString)
	assert.Nil(t, err)

	expiredToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256, JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			},
			UserID: username,
		},
	)

	expiredTokenString, err := expiredToken.SignedString([]byte(config.SecretKey))
	assert.Nil(t, err)

	err = ValidateToken(expiredTokenString)
	assert.Equal(t, ErrorTokenExpired, err)
}

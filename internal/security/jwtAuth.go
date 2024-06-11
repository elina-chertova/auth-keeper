package security

import (
	"errors"
	"github.com/elina-chertova/auth-keeper.git/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const TokenExp = time.Minute * 10

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string
}

var (
	ErrorParseClaims  = errors.New("couldn't parse claims")
	ErrorTokenExpired = errors.New("token expired")
)

func GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
			},
			UserID: username,
		},
	)

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(signedToken string) error {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(
		signedToken, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.SecretKey), nil
		},
	)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return ErrorTokenExpired
			}
		}
		return err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return ErrorParseClaims
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return ErrorTokenExpired
	}
	return nil
}

func GetUserFromToken(signedToken string) (string, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(
		signedToken, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.SecretKey), nil
		},
	)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return "", ErrorParseClaims
	}
	return claims.UserID, nil
}

package security

import (
	"fmt"
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

func TestGenerateToken1(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid username",
			args: args{
				username: "test_user",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := GenerateToken(tt.args.username)
				if !tt.wantErr(t, err, fmt.Sprintf("GenerateToken(%v)", tt.args.username)) {
					return
				}
				assert.NotEmpty(t, got)
			},
		)
	}
}

func TestGetUserFromToken(t *testing.T) {
	username := "test_user"
	validToken, err := GenerateToken(username)
	assert.Nil(t, err)

	type args struct {
		signedToken string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid token",
			args: args{
				signedToken: validToken,
			},
			want:    username,
			wantErr: assert.NoError,
		},
		{
			name: "Invalid token",
			args: args{
				signedToken: "invalid_token",
			},
			want:    "",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := GetUserFromToken(tt.args.signedToken)
				if !tt.wantErr(t, err, fmt.Sprintf("GetUserFromToken(%v)", tt.args.signedToken)) {
					return
				}
				assert.Equalf(t, tt.want, got, "GetUserFromToken(%v)", tt.args.signedToken)
			},
		)
	}
}

func TestValidateToken1(t *testing.T) {
	validUsername := "test_user"
	validToken, err := GenerateToken(validUsername)
	assert.Nil(t, err)

	expiredToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256, JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			},
			UserID: validUsername,
		},
	)
	expiredTokenString, err := expiredToken.SignedString([]byte(config.SecretKey))
	assert.Nil(t, err)

	type args struct {
		signedToken string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid token",
			args: args{
				signedToken: validToken,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Expired token",
			args: args{
				signedToken: expiredTokenString,
			},
			wantErr: assert.Error,
		},
		{
			name: "Invalid token",
			args: args{
				signedToken: "invalid_token",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.wantErr(
					t,
					ValidateToken(tt.args.signedToken),
					fmt.Sprintf("ValidateToken(%v)", tt.args.signedToken),
				)
			},
		)
	}
}

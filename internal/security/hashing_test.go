package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "pass123"
	hash, err := HashPassword(password)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "pass123"
	hash, err := HashPassword(password)
	assert.Nil(t, err)
	assert.True(t, CheckPasswordHash(password, hash))
	assert.False(t, CheckPasswordHash("wrong_password", hash))
}

func TestCheckPasswordHashWithInvalidHash(t *testing.T) {
	password := "pass123"
	invalidHash := "invalid_hash"
	assert.False(t, CheckPasswordHash(password, invalidHash))
}

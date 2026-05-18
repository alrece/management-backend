package auth

import (
	"testing"

	"management-backend/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	config.C = config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret-key-for-unit-testing",
			AccessExpire:  3600,
			RefreshExpire: 86400,
			Issuer:        "management-test",
		},
	}
}

func TestGenerateTokenPair(t *testing.T) {
	pair, err := GenerateTokenPair(1001, "testuser", 1)
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.NotZero(t, pair.AccessExpire)
}

func TestParseAccessToken(t *testing.T) {
	pair, err := GenerateTokenPair(2001, "admin", 2)
	require.NoError(t, err)

	claims, err := ParseAccessToken(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, int64(2001), claims.UserID)
	assert.Equal(t, "admin", claims.Username)
	assert.Equal(t, int64(2), claims.TenantID)
	assert.NotEmpty(t, claims.JTI)
}

func TestParseRefreshToken(t *testing.T) {
	pair, err := GenerateTokenPair(3001, "user", 3)
	require.NoError(t, err)

	claims, err := ParseRefreshToken(pair.RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, int64(3001), claims.UserID)
	assert.Equal(t, "user", claims.Username)
}

func TestParseAccessToken_Invalid(t *testing.T) {
	_, err := ParseAccessToken("invalid.token.string")
	assert.Error(t, err)
}

func TestParseAccessToken_Expired(t *testing.T) {
	origExpire := config.C.JWT.AccessExpire
	config.C.JWT.AccessExpire = -1
	defer func() { config.C.JWT.AccessExpire = origExpire }()

	pair, err := GenerateTokenPair(4001, "expired_user", 1)
	require.NoError(t, err)

	_, err = ParseAccessToken(pair.AccessToken)
	assert.Error(t, err)
}

package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()

	hashedPassword, err := HashPassword("password")
	require.NoError(t, err)
	assert.NotNil(t, hashedPassword)
}

func TestCheckPasswordHash(t *testing.T) {
	t.Parallel()

	hashedPassword, err := HashPassword("password")
	require.NoError(t, err)
	assert.NotNil(t, hashedPassword)

	err = CheckPasswordHash("password", hashedPassword)
	require.NoError(t, err)
	assert.Nil(t, err)
}

func TestMakeJWT(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	require.NoError(t, err)
	assert.NotNil(t, token)
}

func TestValidateJWT(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	require.NoError(t, err)
	assert.NotNil(t, token)

	t.Run("Valid token", func(t *testing.T) {
		t.Parallel()

		userID, err = ValidateJWT(token, tokenSecret)
		require.NoError(t, err)
		assert.NotNil(t, userID)

	})

	t.Run("Invalid token", func(t *testing.T) {
		t.Parallel()

		_, err = ValidateJWT("invalid token", tokenSecret)
		require.Error(t, err)
		assert.NotNil(t, err)

	})

	t.Run("Expired token", func(t *testing.T) {
		t.Parallel()

		expiresIn = -time.Hour
		expiredToken, err := MakeJWT(userID, tokenSecret, expiresIn)
		require.NoError(t, err)

		_, err = ValidateJWT(expiredToken, tokenSecret)
		require.Error(t, err)
		assert.NotNil(t, err)

	})
}

func TestGetBearerToken(t *testing.T) {
	t.Parallel()

	t.Run("Valid token", func(t *testing.T) {
		t.Parallel()
		header := http.Header{}
		header.Add("Authorization", "Bearer token")

		token, err := GetBearerToken(header)
		require.NoError(t, err)
		assert.NotNil(t, token)
	})

	t.Run("Invalid token", func(t *testing.T) {
		t.Parallel()
		header := http.Header{}
		header.Add("Authorization", "Bearer ")

		token, err := GetBearerToken(header)
		require.Error(t, err)
		assert.Equal(t, token, "")
	})
}

func TestRefreshToken(t *testing.T) {
	t.Parallel()

	token, err := MakeRefreshToken()
	require.NoError(t, err)
	assert.NotNil(t, token)
}

func TestGetApiKey(t *testing.T) {
	t.Parallel()

	t.Run("Valid api key", func(t *testing.T) {
		t.Parallel()
		header := http.Header{}
		header.Add("Authorization", "ApiKey apikey")

		apiKey, err := GetAPIKey(header)
		require.NoError(t, err)
		assert.NotNil(t, apiKey)
	})

	t.Run("Invalid api key", func(t *testing.T) {
		t.Parallel()
		header := http.Header{}
		header.Add("Authorization", "ApiKey ")

		apiKey, err := GetAPIKey(header)
		require.Error(t, err)
		assert.Equal(t, apiKey, "")
	})
}

package auth

import (
	"testing"

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

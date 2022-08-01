package hash

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

func TestArgon2Hasher(t *testing.T) {
	hasher := NewArgon2Hasher()
	password := faker.Password()

	firstHashedPassword, err := hasher.Hash(password)
	require.NoError(t, err)
	require.NotEmpty(t, firstHashedPassword)

	require.NoError(t, hasher.Check(password, firstHashedPassword))

	wrongPassword := faker.Password()
	require.EqualError(t, hasher.Check(wrongPassword, firstHashedPassword), ErrIncorrectPassword.Error())

	secondHashedPassword, err := hasher.Hash(password)
	require.NoError(t, err)
	require.NotEmpty(t, secondHashedPassword)
	require.NotEqual(t, firstHashedPassword, secondHashedPassword)
}

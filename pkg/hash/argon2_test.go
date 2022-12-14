package hash_test

import (
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/maypok86/conduit/pkg/hash"
	"github.com/stretchr/testify/require"
)

func TestNewArgon2Hasher(t *testing.T) {
	t.Parallel()

	hasher := hash.NewArgon2Hasher(
		hash.Time(2),
		hash.Memory(32*1024),
		hash.Threads(2),
		hash.KeyLen(16),
		hash.SaltLen(16),
	)
	defaultHasher := hash.NewArgon2Hasher()

	require.False(t, reflect.DeepEqual(hasher, defaultHasher))
}

func TestArgon2Hasher(t *testing.T) {
	t.Parallel()

	hasher := hash.NewArgon2Hasher()
	password := faker.Password()

	firstHashedPassword, err := hasher.Hash(password)
	require.NoError(t, err)
	require.NotEmpty(t, firstHashedPassword)

	require.NoError(t, hasher.Check(password, firstHashedPassword))

	wrongPassword := faker.Password()
	require.EqualError(t, hasher.Check(wrongPassword, firstHashedPassword), hash.ErrIncorrectPassword.Error())

	secondHashedPassword, err := hasher.Hash(password)
	require.NoError(t, err)
	require.NotEmpty(t, secondHashedPassword)
	require.NotEqual(t, firstHashedPassword, secondHashedPassword)
}

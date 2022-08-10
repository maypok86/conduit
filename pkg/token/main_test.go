package token_test

import (
	"testing"
	"time"

	"github.com/maypok86/conduit/pkg/token"
	"github.com/stretchr/testify/require"
)

func fakeToken(t *testing.T, maker token.JWTMaker, email string, duration time.Duration) string {
	t.Helper()

	token, err := maker.CreateToken(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return token
}

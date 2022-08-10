package token_test

import (
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/maypok86/conduit/pkg/token"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	t.Parallel()

	maker, err := token.NewJWTMaker(faker.Password())
	require.NoError(t, err)

	email := faker.Email()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token := fakeToken(t, maker, email, duration)
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	t.Parallel()

	maker, err := token.NewJWTMaker(faker.Password())
	require.NoError(t, err)

	gotToken := fakeToken(t, maker, faker.Email(), -time.Minute)

	payload, err := maker.VerifyToken(gotToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	t.Parallel()

	payload, err := token.NewPayload(faker.Email(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	gotToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := token.NewJWTMaker(faker.Password())
	require.NoError(t, err)

	payload, err = maker.VerifyToken(gotToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrInvalidToken.Error())
	require.Nil(t, payload)
}

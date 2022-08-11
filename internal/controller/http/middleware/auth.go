package middleware

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/internal/controller/http/httperr"
	"github.com/maypok86/conduit/pkg/token"
)

var (
	errAuthHeaderNotProvided   = errors.New("authorization header is not provided")
	errInvalidAuthHeaderFormat = errors.New("invalid authorization header format")
	errPayloadNotFound         = errors.New("payload not found")
)

type tokenMaker interface {
	CreateToken(string, time.Duration) (string, error)
	VerifyToken(token string) (*token.Payload, error)
}

type authMiddleware struct {
	authorizationHeaderKey  string
	authorizationType       string
	authorizationPayloadKey string
	tokenMaker              tokenMaker
}

type authError struct {
	err  error
	slug string
}

func (am authMiddleware) wrapError(err error, slug string) *authError {
	return &authError{
		err:  err,
		slug: slug,
	}
}

func newAuthMiddleware(tokenMaker tokenMaker) authMiddleware {
	return authMiddleware{
		tokenMaker:              tokenMaker,
		authorizationHeaderKey:  "Authorization",
		authorizationType:       "token",
		authorizationPayloadKey: "authorization_payload",
	}
}

func (am authMiddleware) Handle(c *gin.Context) {
	payload, authErr := am.parseAuthHeader(c)
	if authErr != nil {
		httperr.Unauthorised(c, authErr.slug, authErr.err)
	}

	c.Set(am.authorizationPayloadKey, payload)
	c.Next()
}

func (am authMiddleware) getPayload(c *gin.Context) *token.Payload {
	v, exists := c.Get(am.authorizationPayloadKey)
	if exists {
		payload, ok := v.(*token.Payload)
		if ok {
			return payload
		}
	}

	httperr.Unauthorised(c, "not-found-payload", errPayloadNotFound)

	return nil
}

func (am authMiddleware) parseAuthHeader(c *gin.Context) (*token.Payload, *authError) {
	const numberOfFields = 2 // Authorization: Token <token>

	authorizationHeader := c.GetHeader(am.authorizationHeaderKey)
	if authorizationHeader == "" {
		return nil, am.wrapError(errAuthHeaderNotProvided, "empty-token")
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < numberOfFields {
		return nil, am.wrapError(errInvalidAuthHeaderFormat, "invalid-token")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != am.authorizationType {
		return nil, am.wrapError(
			fmt.Errorf("unsupported authorization type %s", authorizationType), //nolint: goerr113
			"unsupported-token",
		)
	}

	accessToken := fields[1]

	payload, err := am.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, am.wrapError(err, "unable-to-verify-jwt")
	}

	return payload, nil
}

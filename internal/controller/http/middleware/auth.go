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

// Auth is a middleware that parses the authorization header and sets the payload to the context.
type Auth struct {
	authorizationHeaderKey  string
	authorizationType       string
	authorizationPayloadKey string
	authorizationTokenKey   string
	tokenMaker              tokenMaker
}

// NewAuth returns a new Auth middleware.
func NewAuth(tokenMaker tokenMaker) Auth {
	return Auth{
		tokenMaker:              tokenMaker,
		authorizationHeaderKey:  "Authorization",
		authorizationType:       "token",
		authorizationPayloadKey: "authorization_payload",
		authorizationTokenKey:   "authorization_token",
	}
}

type authError struct {
	err  error
	slug string
}

func (a Auth) wrapError(err error, slug string) *authError {
	return &authError{
		err:  err,
		slug: slug,
	}
}

// Handle is a middleware that parses the authorization header and sets the payload to the context.
func (a Auth) Handle(c *gin.Context) {
	payload, accessToken, authErr := a.parseAuthHeader(c)
	if authErr != nil {
		httperr.Unauthorised(c, authErr.slug, authErr.err)
		return
	}

	c.Set(a.authorizationTokenKey, accessToken)
	c.Set(a.authorizationPayloadKey, payload)
	c.Next()
}

// GetPayload returns the payload from the context.
func (a Auth) GetPayload(c *gin.Context) *token.Payload {
	v, exists := c.Get(a.authorizationPayloadKey)
	if exists {
		payload, ok := v.(*token.Payload)
		if ok {
			return payload
		}
	}

	httperr.Unauthorised(c, "not-found-payload", errPayloadNotFound)

	return nil
}

// GetToken returns the token from the request.
func (a Auth) GetToken(c *gin.Context) string {
	v, exists := c.Get(a.authorizationTokenKey)
	if exists {
		accessToken, ok := v.(string)
		if ok {
			return accessToken
		}
	}

	return ""
}

func (a Auth) parseAuthHeader(c *gin.Context) (*token.Payload, string, *authError) {
	const numberOfFields = 2 // Authorization: Token <token>

	authorizationHeader := c.GetHeader(a.authorizationHeaderKey)
	if authorizationHeader == "" {
		return nil, "", a.wrapError(errAuthHeaderNotProvided, "empty-token")
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < numberOfFields {
		return nil, "", a.wrapError(errInvalidAuthHeaderFormat, "invalid-token")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != a.authorizationType {
		return nil, "", a.wrapError(
			fmt.Errorf("unsupported authorization type %s", authorizationType), //nolint: goerr113
			"unsupported-token",
		)
	}

	accessToken := fields[1]

	payload, err := a.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, "", a.wrapError(err, "unable-to-verify-jwt")
	}

	return payload, accessToken, nil
}

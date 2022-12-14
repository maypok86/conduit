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
)

//go:generate mockgen -source=auth.go -destination=mocks/auth_test.go -package=middleware_test

// TokenMaker is a token maker.
type TokenMaker interface {
	CreateToken(email string, duration time.Duration) (string, error)
	VerifyToken(accessToken string) (*token.Payload, error)
}

// Auth is a middleware that parses the authorization header and sets the payload to the context.
type Auth struct {
	authorizationHeaderKey  string
	authorizationType       string
	authorizationPayloadKey string
	authorizationTokenKey   string
	tokenMaker              TokenMaker
}

// NewAuth returns a new Auth middleware.
func NewAuth(tokenMaker TokenMaker) Auth {
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

func (a Auth) handle(c *gin.Context, processAuthError func(c *gin.Context, authErr *authError)) {
	payload, accessToken, authErr := a.parseAuthHeader(c)
	if authErr != nil {
		processAuthError(c, authErr)
		return
	}

	c.Set(a.authorizationTokenKey, accessToken)
	c.Set(a.authorizationPayloadKey, payload)
	c.Next()
}

// OptionalHandle is a middleware that optional parses the authorization header and sets the payload to the context.
func (a Auth) OptionalHandle(c *gin.Context) {
	a.handle(c, func(c *gin.Context, authErr *authError) {})
}

// Handle is a middleware that parses the authorization header and sets the payload to the context.
func (a Auth) Handle(c *gin.Context) {
	a.handle(c, func(c *gin.Context, authErr *authError) {
		httperr.Unauthorised(c, authErr.slug, authErr.err)
	})
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

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/pkg/token"
)

type tokenMaker interface {
	CreateToken(string, time.Duration) (string, error)
	VerifyToken(token string) (*token.Payload, error)
}

type authMiddleware struct {
	tokenMaker   tokenMaker
	tokenExpired time.Duration
}

func newAuthMiddleware(tokenMaker tokenMaker, tokenExpired time.Duration) authMiddleware {
	return authMiddleware{
		tokenMaker:   tokenMaker,
		tokenExpired: tokenExpired,
	}
}

func (am authMiddleware) Handle(c *gin.Context) {
	c.Next()
}

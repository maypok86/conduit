// Package httperr provides a http error responses.
package httperr

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/pkg/logger"
	"github.com/maypok86/conduit/pkg/slugerr"
	"go.uber.org/zap"
)

// InternalError is a helper function to respond with internal server error.
func InternalError(c *gin.Context, slug string, err error) {
	httpRespondWithError(c, err, slug, "Internal server error", http.StatusInternalServerError)
}

// Unauthorised is a helper function to respond with unauthorized error.
func Unauthorised(c *gin.Context, slug string, err error) {
	httpRespondWithError(c, err, slug, "Unauthorised", http.StatusUnauthorized)
}

// BadRequest is a helper function to respond with bad request error.
func BadRequest(c *gin.Context, slug string, err error) {
	httpRespondWithError(c, err, slug, "Bad request", http.StatusBadRequest)
}

// RespondWithSlugError is a helper function to respond with slug error.
func RespondWithSlugError(c *gin.Context, err error) {
	var slugError slugerr.SlugError
	if !errors.As(err, &slugError) {
		InternalError(c, "internal-server-error", err)

		return
	}

	switch slugError.ErrorType() {
	case slugerr.ErrorTypeAuthorization:
		Unauthorised(c, slugError.Slug(), slugError)
	case slugerr.ErrorTypeIncorrectInput:
		BadRequest(c, slugError.Slug(), slugError)
	default:
		InternalError(c, slugError.Slug(), slugError)
	}
}

// Errors is a list of http errors.
type Errors struct {
	Body []string `json:"body"`
}

// ErrorResponse is a error response.
type ErrorResponse struct {
	Errors Errors `json:"errors"`
}

func httpRespondWithError(c *gin.Context, err error, slug string, logMessage string, status int) {
	logger.FromRequest(c).Warn(logMessage, zap.Error(err), zap.String("error-slug", slug))
	c.AbortWithStatusJSON(status, ErrorResponse{
		Errors: Errors{
			Body: []string{slug},
		},
	})
}

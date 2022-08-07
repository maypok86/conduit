// Package errors provides custom error types for the openapi HTTP API.
package errors

// ErrorType defines the type of error.
type ErrorType struct {
	t string
}

var (
	// ErrorTypeUnknown defines the unknown type of error.
	ErrorTypeUnknown = ErrorType{"unknown"}
	// ErrorTypeAuthorization defines the authorization type of error.
	ErrorTypeAuthorization = ErrorType{"authorization"}
	// ErrorTypeIncorrectInput defines the incorrect input type of error.
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
)

// SlugError defines error for slug.
type SlugError struct {
	err       string
	slug      string
	errorType ErrorType
}

// Error returns the error message.
func (s SlugError) Error() string {
	return s.err
}

// Slug returns the slug.
func (s SlugError) Slug() string {
	return s.slug
}

// ErrorType returns the error type.
func (s SlugError) ErrorType() ErrorType {
	return s.errorType
}

// NewSlugError creates a new slug error.
func NewSlugError(err string, slug string) SlugError {
	return SlugError{
		err:       err,
		slug:      slug,
		errorType: ErrorTypeUnknown,
	}
}

// NewAuthorizationError creates a new authorization error.
func NewAuthorizationError(err string, slug string) SlugError {
	return SlugError{
		err:       err,
		slug:      slug,
		errorType: ErrorTypeAuthorization,
	}
}

// NewIncorrectInputError creates a new incorrect input error.
func NewIncorrectInputError(err string, slug string) SlugError {
	return SlugError{
		err:       err,
		slug:      slug,
		errorType: ErrorTypeIncorrectInput,
	}
}

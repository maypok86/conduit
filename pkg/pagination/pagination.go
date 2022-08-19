// Package pagination provides pagination.
package pagination

const (
	// DefaultLimit is a default limit for pagination.
	DefaultLimit = 20
	// MaxLimit is a max limit for pagination.
	MaxLimit = 100
)

// Params is a params for pagination.
type Params struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

// List is a list of pagination.
type List[T any] struct {
	Result  []T  `json:"result"`
	HasNext bool `json:"has_next"`
}

// NewList returns a new list of pagination.
func NewList[T any](objectList []T, limit uint64) List[T] {
	objectLength := uint64(len(objectList))

	hasNext := objectLength == limit+1
	if hasNext {
		objectList = objectList[:objectLength-1]
	}

	return List[T]{
		Result:  objectList,
		HasNext: hasNext,
	}
}

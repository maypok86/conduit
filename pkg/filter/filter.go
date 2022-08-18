// Package filter provides filters for squirrel.
package filter

import (
	sq "github.com/Masterminds/squirrel"
)

// Type is a type of filter.
type Type uint8

const (
	// TypeEQ value is equal.
	TypeEQ Type = iota + 1

	// TypeNotEQ value is not equal.
	TypeNotEQ

	// TypeGTE greater than value.
	TypeGTE

	// TypeGT value greater.
	TypeGT

	// TypeLT value less.
	TypeLT

	// TypeLTE value less or equals.
	TypeLTE

	// TypeLike value can contain.
	TypeLike

	// TypeNotLike value cannot contain.
	TypeNotLike

	// TypeILike value can contain (case insensitive).
	TypeILike

	// TypeNotILike value cannot contain (case insensitive).
	TypeNotILike
)

// Operator is an operator for linking filters.
type Operator uint8

const (
	// OperatorAnd is an operator for AND linking filters.
	OperatorAnd Operator = iota + 1
	// OperatorOr is an operator for OR linking filters.
	OperatorOr
)

// Filter is a filter for squirrel.SelectBuilder.
type Filter struct {
	column   string
	ftype    Type
	operator Operator
	filters  []Filter
	value    any
}

// New creates a new filter.
func New(column string, ftype Type, value any) Filter {
	return Filter{
		column:   column,
		ftype:    ftype,
		value:    value,
		operator: OperatorAnd,
		filters:  make([]Filter, 0),
	}
}

// SetOperator sets operator for linking filters.
func (f Filter) SetOperator(operator Operator) Filter {
	f.operator = operator

	return f
}

// WithFilters adds filters.
func (f Filter) WithFilters(filters ...Filter) Filter {
	f.filters = append(f.filters, filters...)

	return f
}

func (f Filter) condition() sq.Sqlizer { //nolint:cyclop,ireturn
	switch f.ftype {
	case TypeEQ:
		return sq.Eq{f.column: f.value}
	case TypeNotEQ:
		return sq.NotEq{f.column: f.value}
	case TypeGTE:
		return sq.GtOrEq{f.column: f.value}
	case TypeGT:
		return sq.Gt{f.column: f.value}
	case TypeLT:
		return sq.Lt{f.column: f.value}
	case TypeLTE:
		return sq.LtOrEq{f.column: f.value}
	case TypeLike:
		return sq.Like{f.column: f.value}
	case TypeNotLike:
		return sq.NotLike{f.column: f.value}
	case TypeILike:
		return sq.ILike{f.column: f.value}
	case TypeNotILike:
		return sq.NotILike{f.column: f.value}
	}

	return sq.Eq{f.column: f.value}
}

func (f Filter) getConditions() sq.Sqlizer { //nolint:ireturn
	condition := f.condition()

	if len(f.filters) == 0 {
		return condition
	}

	conditions := []sq.Sqlizer{condition}

	for _, filter := range f.filters {
		conditions = append(conditions, filter.getConditions())
	}

	if f.operator == OperatorOr {
		return or(conditions)
	}

	return and(conditions)
}

// UseSelectBuilder adds filters to squirrel.SelectBuilder.
func (f Filter) UseSelectBuilder(sb sq.SelectBuilder) sq.SelectBuilder {
	return sb.Where(f.getConditions())
}

func and(conditions []sq.Sqlizer) sq.And {
	var result sq.And

	for _, condition := range conditions {
		result = append(result, condition)
	}

	return result
}

func or(conditions []sq.Sqlizer) sq.Or {
	var result sq.Or

	for _, condition := range conditions {
		result = append(result, condition)
	}

	return result
}

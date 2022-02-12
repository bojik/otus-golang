package hw09structvalidator

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidType                 = errors.New("invalid type of variable, expected struct")
	ErrExpectedAtLeastOneParameter = errors.New("expected at least one parameter")
	ErrNotAppropriatedType         = errors.New("validator is not appropriate for this type")
)

type ParseValidatorError struct {
	Tag string
	Err error
}

func (p ParseValidatorError) Error() string {
	return fmt.Sprintf("invalid validator `%s`: %s", p.Tag, p.Err.Error())
}

var _ error = (*ParseValidatorError)(nil)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

var _ error = (*ValidationErrors)(nil)
var _ error = (*ValidationError)(nil)

func NewParseValidatorError(err error) *ParseValidatorError {
	return &ParseValidatorError{Err: err}
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{Err: err}
}

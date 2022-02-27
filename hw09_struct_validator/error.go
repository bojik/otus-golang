package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

var (
	ErrExpectedOneParameter        = errors.New("expected one parameter")
	ErrValidatorDoesNotExist       = errors.New("validator does not exist")
	ErrInvalidParameter            = errors.New("invalid parameter of validator")
	ErrExpectedAtLeastOneParameter = errors.New("expected at least one parameter")
	ErrInvalidType                 = errors.New("invalid type of variable, expected struct")
	ErrInappropriateType           = errors.New("validator is not appropriate for this type")
)

var (
	ErrValidationIntRange    = errors.New("invalid int range")
	ErrValidationStringLen   = errors.New("invalid string len")
	ErrValidationMin         = errors.New("min validator error")
	ErrValidationMax         = errors.New("max validator error")
	ErrValidationStringRange = errors.New("invalid string range")
	ErrValidationRegexp      = errors.New("string does not match regexp")
)

type ParseValidatorError struct {
	Tag string
	Err error
}

func (p ParseValidatorError) Unwrap() error {
	return p.Err
}

func (p ParseValidatorError) Error() string {
	return fmt.Sprintf("invalid validator `%s`: %s", p.Tag, p.Err.Error())
}

var (
	_ error           = (*ParseValidatorError)(nil)
	_ xerrors.Wrapper = (*ParseValidatorError)(nil)
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Unwrap() error {
	return v.Err
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err)
}

var (
	_ error           = (*ValidationError)(nil)
	_ xerrors.Wrapper = (*ValidationError)(nil)
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	msgs := []string{}
	for _, e := range v {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "\n")
}

var _ error = (*ValidationErrors)(nil)

func NewParseValidatorError(err error) ParseValidatorError {
	return ParseValidatorError{Err: err}
}

func NewValidationError(err error) ValidationError {
	return ValidationError{Err: err}
}

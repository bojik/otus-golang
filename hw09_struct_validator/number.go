package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type MinValidator struct {
	ErrMessage string
}

func (m MinValidator) Validate(v interface{}, args ...string) error {
	msg := m.ErrMessage
	if msg == "" {
		msg = "value %s is less then %s"
	}
	if len(args) < 1 {
		return NewParseValidatorError(ErrExpectedAtLeastOneParameter)
	}
	min, err := strconv.Atoi(args[0])
	if err != nil {
		return NewParseValidatorError(err)
	}
	r := reflect.ValueOf(v)
	switch r.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if r.Int() < int64(min) {
			return NewValidationError(errors.New(fmt.Sprintf(msg, r.Int(), min)))
		}
	case reflect.Float32, reflect.Float64:
		if r.Float() < float64(min) {
			return NewValidationError(errors.New(fmt.Sprintf(msg, r.Float(), min)))
		}
	default:
		return NewParseValidatorError(ErrNotAppropriatedType)
	}
	return nil
}

var _ Validator = (*MinValidator)(nil)

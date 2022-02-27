package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type (
	RegexpValidation      struct{}
	StringRangeValidation struct{}
	StringLenValidation   struct{}
)

func (validator RegexpValidation) Validate(v interface{}, args ...string) error {
	if len(args) != 1 {
		return NewParseValidatorError(ErrExpectedOneParameter)
	}
	re, err := regexp.Compile(args[0])
	if err != nil {
		return NewParseValidatorError(fmt.Errorf("%w: %s", ErrInvalidParameter, err.Error()))
	}
	refValue := reflect.ValueOf(v)
	//nolint: exhaustive
	switch refValue.Kind() {
	case reflect.String:
		if !re.MatchString(refValue.String()) {
			return NewValidationError(
				fmt.Errorf("%w: '%s'", ErrValidationRegexp, re.String()),
			)
		}
	case reflect.Slice:
		for i := 0; i < refValue.Len(); i++ {
			item := refValue.Index(i)
			if !item.CanInterface() {
				return NewValidationError(fmt.Errorf("%w: %s", ErrInappropriateType, item.Type()))
			}
			if err := validator.Validate(item.Interface(), args...); err != nil {
				return err
			}
		}
	default:
		return NewParseValidatorError(ErrInappropriateType)
	}
	return nil
}

func (validator StringRangeValidation) Validate(v interface{}, args ...string) error {
	if len(args) < 1 {
		return NewParseValidatorError(ErrExpectedAtLeastOneParameter)
	}
	refValue := reflect.ValueOf(v)
	//nolint: exhaustive
	switch refValue.Kind() {
	case reflect.String:
		for _, s := range args {
			if refValue.String() == s {
				return nil
			}
		}
		return NewValidationError(
			fmt.Errorf(
				"%w: %s not in %s",
				ErrValidationStringRange,
				refValue.String(),
				strings.Join(args, ", "),
			),
		)
	case reflect.Slice:
		for i := 0; i < refValue.Len(); i++ {
			item := refValue.Index(i)
			if !item.CanInterface() {
				return NewValidationError(fmt.Errorf("%w: %s", ErrInappropriateType, item.Type()))
			}
			if err := validator.Validate(item.Interface(), args...); err != nil {
				return err
			}
		}
	default:
		return NewParseValidatorError(ErrInappropriateType)
	}
	return nil
}

func (validator StringLenValidation) Validate(v interface{}, args ...string) error {
	if len(args) != 1 {
		return NewParseValidatorError(ErrExpectedOneParameter)
	}
	length, err := strconv.Atoi(args[0])
	if err != nil {
		return NewParseValidatorError(fmt.Errorf("%w: %s", ErrInvalidParameter, err.Error()))
	}
	if length < 0 {
		return NewParseValidatorError(fmt.Errorf("%w: length < 0", ErrInvalidParameter))
	}
	refValue := reflect.ValueOf(v)
	//nolint: exhaustive
	switch refValue.Kind() {
	case reflect.String:
		if refValue.Len() != length {
			return NewValidationError(
				fmt.Errorf(
					"%w: '%s' has invalid lengh %d, expected %d",
					ErrValidationStringLen,
					refValue.String(),
					refValue.Len(),
					length,
				),
			)
		}
	case reflect.Slice:
		for i := 0; i < refValue.Len(); i++ {
			item := refValue.Index(i)
			if !item.CanInterface() {
				return NewValidationError(fmt.Errorf("%w: %s", ErrInappropriateType, item.Type()))
			}
			if err := validator.Validate(item.Interface(), args...); err != nil {
				return err
			}
		}
	default:
		return NewParseValidatorError(ErrInappropriateType)
	}
	return nil
}

var (
	_ Validator = (*RegexpValidation)(nil)
	_ Validator = (*StringLenValidation)(nil)
	_ Validator = (*StringRangeValidation)(nil)
)

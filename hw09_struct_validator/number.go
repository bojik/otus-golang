package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type (
	MinValidation      struct{}
	MaxValidation      struct{}
	IntRangeValidation struct{}
)

func (validator MinValidation) Validate(v interface{}, args ...string) error {
	return validateMinMax(
		validator,
		v,
		args,
		func(v, min int64) bool {
			return v >= min
		},
		func(v, min int64) error {
			return fmt.Errorf("%w: %d < %d", ErrValidationMin, v, min)
		},
	)
}

func (validator MaxValidation) Validate(v interface{}, args ...string) error {
	return validateMinMax(
		validator,
		v,
		args,
		func(v, max int64) bool {
			return v <= max
		},
		func(v, max int64) error {
			return fmt.Errorf("%w: %d > %d", ErrValidationMax, v, max)
		},
	)
}

func validateMinMax(
	validator Validator,
	v interface{},
	args []string,
	isValid func(v, max int64) bool,
	errorFormatter func(v, max int64) error,
) error {
	if len(args) != 1 {
		return NewParseValidatorError(ErrExpectedOneParameter)
	}
	mm, err := strconv.Atoi(args[0])
	if err != nil {
		return NewParseValidatorError(fmt.Errorf("%w: %s", ErrInvalidParameter, err.Error()))
	}
	refValue := reflect.ValueOf(v)
	//nolint: exhaustive
	switch refValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !isValid(refValue.Int(), int64(mm)) {
			return NewValidationError(errorFormatter(refValue.Int(), int64(mm)))
		}
	case reflect.Slice:
		for i := 0; i < refValue.Len(); i++ {
			item := refValue.Index(i)
			if !item.CanInterface() {
				return fmt.Errorf("%w: %s", ErrInappropriateType, item.Type())
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

func (validator IntRangeValidation) Validate(v interface{}, args ...string) error {
	if len(args) < 1 {
		return NewParseValidatorError(ErrExpectedAtLeastOneParameter)
	}
	arr := []int64{}
	for _, arg := range args {
		n, err := strconv.Atoi(arg)
		if err != nil {
			return NewParseValidatorError(fmt.Errorf("%w: %s", ErrInvalidParameter, err.Error()))
		}
		arr = append(arr, int64(n))
	}
	refValue := reflect.ValueOf(v)
	//nolint: exhaustive
	switch refValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for _, val := range arr {
			if val == refValue.Int() {
				return nil
			}
		}
		return NewValidationError(
			fmt.Errorf(
				"%w: %d not in range (%s)",
				ErrValidationIntRange,
				refValue.Int(),
				strings.Join(args, ", "),
			),
		)
	case reflect.Slice:
		for i := 0; i < refValue.Len(); i++ {
			item := refValue.Index(i)
			if !item.CanInterface() {
				return fmt.Errorf("%w: %s", ErrInappropriateType, item.Type())
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
	_ Validator = (*MinValidation)(nil)
	_ Validator = (*MaxValidation)(nil)
	_ Validator = (*IntRangeValidation)(nil)
)

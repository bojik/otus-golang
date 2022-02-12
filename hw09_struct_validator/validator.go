package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type Validator interface {
	Validate(v interface{}, args ...string) error
}

var mapping = map[string]Validator{
	"in":     InValidation{},
	"min":    MinValidation{},
	"max":    MaxValidation{},
	"regexp": RegexpValidation{},
	"len":    StringLenValidation{},
}

func Validate(v interface{}) error {
	r := reflect.TypeOf(v)
	if r.Kind() != reflect.Struct {
		return ErrInvalidType
	}
	rv := reflect.ValueOf(v)
	var errs ValidationErrors
	for i := 0; i < r.NumField(); i++ {
		v, ok := r.Field(i).Tag.Lookup("validate")
		if !ok || !rv.Field(i).CanInterface() {
			continue
		}
		tags := ParseTag(v)
		for _, tag := range tags {
			err := validateByTag(tag, rv.Field(i))
			switch e := err.(type) { // nolint:errorlint
			case ValidationError:
				e.Field = r.Field(i).Name
				errs = append(errs, e) // collect error of validation
			case ParseValidatorError:
				e.Tag = tag.OriginalTag
				return e // it seems, someone needs to fix validator declaration
			}
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func validateByTag(tag TagData, val reflect.Value) error {
	validator, ok := mapping[tag.Name]
	if !ok {
		return NewParseValidatorError(fmt.Errorf("%w: %s", ErrValidatorDoesNotExist, tag.Name))
	}
	return validator.Validate(val.Interface(), tag.Args...)
}

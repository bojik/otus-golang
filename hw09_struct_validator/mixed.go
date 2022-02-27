package hw09structvalidator

import "reflect"

type InValidation struct{}

func (InValidation) Validate(v interface{}, args ...string) error {
	refType := reflect.TypeOf(v)
	var validator Validator
	if refType.Kind() == reflect.String {
		validator = StringRangeValidation{}
	} else {
		validator = IntRangeValidation{}
	}
	return validator.Validate(v, args...)
}

var _ Validator = (*InValidation)(nil)

package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type Validator interface {
	Validate(v interface{}, args ...string) error
}

func Validate(v interface{}) error {
	r := reflect.TypeOf(v)
	if r.Kind() != reflect.Struct {
		return ErrInvalidType
	}
	rv := reflect.ValueOf(v)
	for i := 0; i < r.NumField(); i++ {
		v, ok := r.Field(i).Tag.Lookup("validate")
		if ok {
			fmt.Println(r.Field(i).Name, v, rv.Field(i).String())
		}
	}
	return nil
}

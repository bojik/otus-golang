package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegexpValidation_Validate(t *testing.T) {
	tests := []struct {
		v    interface{}
		args []string
		err  error
	}{
		{"123", []string{}, ErrExpectedOneParameter},
		{"123", []string{"\\"}, ErrInvalidParameter},
		{"123a", []string{"^\\d+$"}, ErrValidationRegexp},
		{123, []string{"^\\d+$"}, ErrInappropriateType},
		{[]string{"123", "456", "123z"}, []string{"^\\d+$"}, ErrValidationRegexp},
		{"123", []string{"\\d+"}, nil},
		{[]string{"123", "456", "123"}, []string{"^\\d+$"}, nil},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("min: %v args: %v", tc.v, tc.args), func(t *testing.T) {
			v := RegexpValidation{}
			err := v.Validate(tc.v, tc.args...)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestStringRangeValidation_Validate(t *testing.T) {
	tests := []struct {
		v    interface{}
		args []string
		err  error
	}{
		{"123", []string{}, ErrExpectedAtLeastOneParameter},
		{"123", []string{"1", "2", "3"}, ErrValidationStringRange},
		{[]interface{}{"123", "456", 123}, []string{"789", "456", "123"}, ErrInappropriateType},
		{"123", []string{"789", "456", "123"}, nil},
		{[]string{"123", "456", "789"}, []string{"789", "456", "123"}, nil},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("min: %v args: %v", tc.v, tc.args), func(t *testing.T) {
			v := StringRangeValidation{}
			err := v.Validate(tc.v, tc.args...)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestStringLenValidation_Validate(t *testing.T) {
	tests := []struct {
		v    interface{}
		args []string
		err  error
	}{
		{"123", []string{}, ErrExpectedOneParameter},
		{"123", []string{"1", "2", "3"}, ErrExpectedOneParameter},
		{"123", []string{"-1"}, ErrInvalidParameter},
		{"123", []string{"abc"}, ErrInvalidParameter},
		{[]interface{}{"123", "456", "789", 567}, []string{"3"}, ErrInappropriateType},
		{"123", []string{"2"}, ErrValidationStringLen},
		{"123", []string{"3"}, nil},
		{[]interface{}{"123", "456", "789"}, []string{"3"}, nil},
		{[]interface{}{"123", "45", "789"}, []string{"3"}, ErrValidationStringLen},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("min: %v args: %v", tc.v, tc.args), func(t *testing.T) {
			v := StringLenValidation{}
			err := v.Validate(tc.v, tc.args...)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

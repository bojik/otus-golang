package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMinValidator_Validate(t *testing.T) {
	tests := []struct {
		v    interface{}
		args []string
		err  error
	}{
		{12, []string{"15"}, ErrValidationMin},
		{16, []string{"abs"}, ErrInvalidParameter},
		{16.54, []string{"15"}, ErrInappropriateType},
		{16, []string{}, ErrExpectedOneParameter},
		{[]int{44, 22, 33}, []string{"23"}, ErrValidationMin},
		{[]interface{}{"test", 123}, []string{"1"}, ErrInappropriateType},
		{[]interface{}{struct{ test string }{"123"}}, []string{"1"}, ErrInappropriateType},
		{16, []string{"15"}, nil},
		{0, []string{"-1"}, nil},
		{[]int{44, 22, 33}, []string{"21"}, nil},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("min: %v args: %v", tc.v, tc.args), func(t *testing.T) {
			v := MinValidation{}
			err := v.Validate(tc.v, tc.args...)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestMaxValidator_Validate(t *testing.T) {
	tests := []struct {
		v    interface{}
		args []string
		err  error
	}{
		{12, []string{"11"}, ErrValidationMax},
		{16, []string{}, ErrExpectedOneParameter},
		{16, []string{"abs"}, ErrInvalidParameter},
		{16.54, []string{"15"}, ErrInappropriateType},
		{[]int{44, 22, 33}, []string{"43"}, ErrValidationMax},
		{[]interface{}{"test", 123}, []string{"1"}, ErrInappropriateType},
		{[]interface{}{struct{ test string }{"123"}}, []string{"1"}, ErrInappropriateType},
		{16, []string{"17"}, nil},
		{-2, []string{"-1"}, nil},
		{[]int{44, 22, 33}, []string{"45"}, nil},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("max: %v args: %v", tc.v, tc.args), func(t *testing.T) {
			v := MaxValidation{}
			err := v.Validate(tc.v, tc.args...)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestIntRangeValidator_Validate(t *testing.T) {
	tests := []struct {
		v    interface{}
		args []string
		err  error
	}{
		{16, []string{}, ErrExpectedAtLeastOneParameter},
		{12, []string{"10", "11"}, ErrValidationIntRange},
		{9, []string{"10", "11"}, ErrValidationIntRange},
		{16, []string{"abs", "12"}, ErrInvalidParameter},
		{16, []string{"12", "abs"}, ErrInvalidParameter},
		{16.54, []string{"15", "17"}, ErrInappropriateType},
		{[]int{44, 22, 33}, []string{"20", "43"}, ErrValidationIntRange},
		{[]interface{}{"test", 123}, []string{"1", "100"}, ErrInappropriateType},
		{[]interface{}{struct{ test string }{"123"}}, []string{"1", "2"}, ErrInappropriateType},
		{16, []string{"16", "17"}, nil},
		{-2, []string{"-3", "-2", "-1"}, nil},
		{[]int{44, 22, 33}, []string{"22", "44", "33"}, nil},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("range: %v args: %v", tc.v, tc.args), func(t *testing.T) {
			v := IntRangeValidation{}
			err := v.Validate(tc.v, tc.args...)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

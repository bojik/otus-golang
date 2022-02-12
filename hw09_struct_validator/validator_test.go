package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{App{"12345"}, nil},
		{
			App{"1234"},
			ValidationErrors{
				ValidationError{Err: ErrValidationStringLen, Field: "Version"},
			},
		},
		{
			User{
				ID:     "2030edfd-1610-44d1-a7fd-c08efab88a3e",
				Name:   "Test",
				Age:    18,
				Email:  "test@test.com",
				Role:   UserRole("admin"),
				Phones: []string{"12345678901"},
				meta:   nil,
			},
			nil,
		},
		{
			User{
				ID:     "2030edfd-1610-44d1-a7fd-c08efab88a3e",
				Name:   "Test",
				Age:    51,
				Email:  "test@test.com",
				Role:   UserRole("admin"),
				Phones: []string{"12345678901"},
			},
			ValidationErrors{
				ValidationError{Err: ErrValidationMax, Field: "Age"},
			},
		},
		{
			User{
				ID:     "2030edfd-1610-44d1-a7fd-c08efab88a3",
				Name:   "Test",
				Age:    15,
				Email:  "test.com",
				Role:   UserRole("admin1"),
				Phones: []string{"12345678901", "1234567890"},
			},
			ValidationErrors{
				ValidationError{Err: ErrValidationStringLen, Field: "ID"},
				ValidationError{Err: ErrValidationMin, Field: "Age"},
				ValidationError{Err: ErrValidationRegexp, Field: "Email"},
				ValidationError{Err: ErrValidationStringRange, Field: "Role"},
				ValidationError{Err: ErrValidationStringLen, Field: "Phones"},
			},
		},
		{
			Response{Code: 201, Body: "unk"}, ValidationErrors{ValidationError{Err: ErrValidationIntRange, Field: "Code"}},
		},
		{Token{}, nil},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.Nil(t, err)
				return
			}
			require.IsType(t, ValidationErrors{}, err)

			expectedErr, ok := tt.expectedErr.(ValidationErrors) // nolint:errorlint
			require.True(t, ok)

			originalErr, ok := err.(ValidationErrors) // nolint:errorlint
			require.True(t, ok)
			require.Len(t, originalErr, len(expectedErr))

			for i := 0; i < len(expectedErr); i++ {
				exp := expectedErr[i]
				orig := originalErr[i]
				require.Equal(t, exp.Field, orig.Field)
				require.ErrorIs(t, orig.Err, exp.Err)
			}
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		data        interface{}
		expectedErr error
		tag         string
	}{
		{struct {
			Age int `validate:"min"`
		}{18}, ErrExpectedOneParameter, "min"},
		{struct {
			Age int `validate:"in:"`
		}{18}, ErrExpectedAtLeastOneParameter, "in:"},
		{struct {
			Age int `validate:"in2:"`
		}{18}, ErrValidatorDoesNotExist, "in2:"},
		{struct {
			Age int `validate:"regexp:\\"`
		}{18}, ErrInvalidParameter, "regexp:\\"},
		{struct {
			Age int `validate:"regexp:\\d+"`
		}{18}, ErrInappropriateType, "regexp:\\d+"},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("test case: #%d", i), func(t *testing.T) {
			err := Validate(tc.data)
			actErr, ok := err.(ParseValidatorError) // nolint:errorlint
			require.True(t, ok)
			require.ErrorIs(t, err, tc.expectedErr)
			if tc.tag != "" {
				require.Equal(t, tc.tag, actErr.Tag)
			}
		})
	}
	require.ErrorIs(t, Validate("123"), ErrInvalidType)
}

package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInValidation_Validate(t *testing.T) {
	tests := []struct {
		v    interface{}
		args []string
		err  error
	}{
		{12, []string{"13", "15"}, ErrValidationIntRange},
		{"12", []string{"13", "15"}, ErrValidationStringRange},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("min: %v args: %v", tc.v, tc.args), func(t *testing.T) {
			v := InValidation{}
			err := v.Validate(tc.v, tc.args...)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

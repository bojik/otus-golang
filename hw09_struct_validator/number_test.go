package hw09structvalidator

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMinValidator_Validate(t *testing.T) {
	v := MinValidator{}
	err := v.Validate(12, "15")
	//	var target *ValidationError
	//fmt.Println(err)
	require.ErrorAs(t, err, &ValidationError{})
}

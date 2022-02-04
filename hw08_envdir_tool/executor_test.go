package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	exp := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: " ", NeedRemove: false},
		"FOO":   EnvValue{Value: "   foo\nwith", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{Value: "", NeedRemove: true},
	}
	os.Setenv("HELLO", "hzhz")
	os.Setenv("EMPTY", "empty")
	temp, err := os.CreateTemp("./testdata", "temp.*")
	require.Nil(t, err)
	defer os.Remove(temp.Name())
	stdout := os.Stdout
	os.Stdout = temp
	code := RunCmd([]string{"./testdata/echo.sh", "test1", "test2"}, exp)
	require.Equal(t, 0, code)
	temp.Seek(0, io.SeekStart)
	data, _ := io.ReadAll(temp)
	temp.Close()
	os.Stdout = stdout
	require.Equal(
		t,
		"HELLO is (\"hello\")\nBAR is (bar)\nFOO is (   foo\nwith)\nUNSET is ()\n"+
			"ADDED is ()\nEMPTY is ( )\narguments are test1 test2\n",
		string(data),
	)
}

func TestErrorRunCmd(t *testing.T) {
	code := RunCmd([]string{"ls", "hzhz"}, make(Environment))
	require.NotEqual(t, 0, code)
}

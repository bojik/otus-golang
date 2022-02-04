package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	testDir := "./testdata/env"
	invFile := testDir + string(os.PathSeparator) + "HZB=HZB"
	err := os.WriteFile(invFile, []byte("test"), 0x666)
	require.Nil(t, err)
	defer os.Remove(invFile)
	res, err1 := ReadDir(testDir)
	require.Nil(t, err1)
	_, ok := res["HZB=HZB"]
	require.False(t, ok)
	exp := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: " ", NeedRemove: false},
		"FOO":   EnvValue{Value: "   foo\nwith", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{Value: "", NeedRemove: true},
		//		"UNSET1": EnvValue{Value: "", NeedRemove: true},
	}
	require.Equal(t, len(exp), len(res))
	for k, v := range exp {
		v1, ok := res[k]
		if !ok {
			t.Errorf("key not found: %s", k)
			t.FailNow()
		}
		require.Equal(t, v.Value, v1.Value)
		require.Equal(t, v.NeedRemove, v1.NeedRemove)
	}
}

func TestReadDirError(t *testing.T) {
	t.Run("err dir does not exist", func(t *testing.T) {
		res, err := ReadDir("hzhz")
		require.Nil(t, res)
		require.ErrorIs(t, err, ErrDirDoesNotExist)
	})
	t.Run("err is not dir", func(t *testing.T) {
		res, err := ReadDir("./go.mod")
		require.Nil(t, res)
		require.ErrorIs(t, err, ErrIsNotDir)
	})
}

func TestClearText(t *testing.T) {
	cases := []struct {
		s string
		e string
	}{
		{"test", "test"},
		{"test\t ", "test"},
		{"test\n \t ", "test\n"},
		{"\t  super test\t \n", "\t  super test\t \n"},
		{"\t  super test\n\t  ", "\t  super test\n"},
		{"\t  super test\n\x00\t  ", "\t  super test\n\n"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run("clear text "+tc.s, func(t *testing.T) {
			require.Equal(t, tc.e, clearString(tc.s))
		})
	}
}

func TestReadFirstLine(t *testing.T) {
	fl, err := readFirstLine("./testdata/env/FOO")
	require.Nil(t, err)
	require.Equal(t, "   foo\nwith", clearString(fl))
}

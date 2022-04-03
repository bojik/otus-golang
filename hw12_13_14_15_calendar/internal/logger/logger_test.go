package logger

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestLogger(t *testing.T) {
	f := "./testdata/app.log"
	logger := New(
		&OptionMinLevel{INFO},
	)
	logger.ResetWriters()
	fp, err := logger.AddLogFile(f)
	defer func() {
		fp.Close()
		os.Remove(f)
	}()
	require.Nil(t, err)
	lg := []struct {
		lvl            Level
		msg            string
		params         []Parameter
		expectedParams string
	}{
		{
			DEBUG, "Debug message",
			[]Parameter{NewIntParam("test", 123)},
			"[test=123]",
		},
		{
			INFO,
			"Info message",
			[]Parameter{NewInt32Param("test", 123), NewStringParam("s", "test_str")},
			"[test=123, s=test_str]",
		},
		{
			ERROR,
			"Error message",
			[]Parameter{NewFloat32Param("float", 2.2)},
			"[float=2.200000]",
		},
	}
	logger.Debug(lg[0].msg, lg[0].params...)
	logger.Info(lg[1].msg, lg[1].params...)
	logger.Error(lg[2].msg, lg[2].params...)
	content, err := ioutil.ReadFile(f)
	lines := strings.Split(string(content), "\n")
	require.Len(t, lines, len(lg))
	for i, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		require.Equal(t, lg[i+1].lvl.String(), parts[1])
		require.Equal(t, lg[i+1].msg, parts[2])
		require.Equal(t, lg[i+1].expectedParams, parts[3])
	}
}

func TestLoggerMultiThreads(t *testing.T) {
	f := "./testdata/app.log"
	logger := New(
		&OptionMinLevel{INFO},
	)
	logger.ResetWriters()
	fp, err := logger.AddLogFile(f)
	require.NoError(t, err)
	defer func() {
		fp.Close()
		os.Remove(f)
	}()
	c := 100
	wg := sync.WaitGroup{}
	wg.Add(c)
	for i := 0; i < c; i++ {
		i := i
		go func() {
			defer wg.Done()
			logger.Info(fmt.Sprintf("Thread %d", i))
		}()
	}
	wg.Wait()
	content, err := ioutil.ReadFile(f)
	lines := strings.Split(strings.TrimRight(string(content), "\n"), "\n")
	require.Equal(t, c, len(lines))
}

func TestLoggerMessage(t *testing.T) {
	f := "./testdata/app.log"
	logger := New(
		&OptionMinLevel{DEBUG},
	)
	logger.ResetWriters()
	fp, err := logger.AddLogFile(f)
	defer func() {
		fp.Close()
		os.Remove(f)
	}()
	require.Nil(t, err)
	lg := []struct {
		lvl Level
		msg string
	}{
		{DEBUG, "Debug message"},
		{INFO, "Info message"},
		{ERROR, "Error message"},
	}
	for _, l := range lg {
		logger.Save(l.lvl, l.msg)
	}
	content, err := ioutil.ReadFile(f)
	lines := strings.Split(string(content), "\n")
	require.Len(t, lines, len(lg)+1)
	for i, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		require.Equal(t, lg[i].lvl.String(), parts[1])
		require.Equal(t, lg[i].msg, parts[2])
	}
}

func TestLoggerOptions(t *testing.T) {
	logger := New(
		&OptionMinLevel{DEBUG},
	)
	require.Equal(t, DEBUG, logger.(*logg).minLevel)
}

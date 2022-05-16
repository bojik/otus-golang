//go:build integration
// +build integration

package tests

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	EnvHTTPUrl        = "CALENDAR_HTTP_URL"
	EnvAPIUrl         = "CALENDAR_API_URL"
	EnvSelectInterval = "CALENDAR_SELECT_INTERVAL"
)

const (
	FormatDateTime = "2006-01-02 15:04:05"
	FormatDate     = "2006-01-02"
)

type response struct {
	Error   string
	Success bool
	Data    struct {
		Event struct {
			ID    string
			Title string
			Sent  bool
		}
	}
}

type listResponse struct {
	Error   string
	Success bool
	Data    struct {
		Events []struct {
			ID    string
			Title string
		}
	}
}

type suiteType interface {
	T() *testing.T
}

func getHTTPUrl(s suiteType) string {
	url, ok := os.LookupEnv(EnvHTTPUrl)
	require.Truef(s.T(), ok, "Environment %s is not set", EnvHTTPUrl)
	//	url := "http://0.0.0.0:8085"
	return url
}

func getApiUrl(s suiteType) string {
	url, ok := os.LookupEnv(EnvAPIUrl)
	require.Truef(s.T(), ok, "Environment %s is not set", EnvAPIUrl)
	//	url := "0.0.0.0:8086"
	return url
}

func getSelectInterval(s suiteType) time.Duration {
	interval, ok := os.LookupEnv(EnvSelectInterval)
	require.Truef(s.T(), ok, "Environment %s is not set", EnvSelectInterval)
	//	interval := "2s"
	d, err := time.ParseDuration(interval)
	require.Nil(s.T(), err)
	return d
}

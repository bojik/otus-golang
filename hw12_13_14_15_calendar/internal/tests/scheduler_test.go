//go:build integration
// +build integration

package tests

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SchedulerTestSuite struct {
	suite.Suite
	url string
}

func (s *SchedulerTestSuite) SetupSuite() {
	s.url = getHTTPUrl(s)
}

func (s *SchedulerTestSuite) TestScheduler() {
	t1, _ := time.Parse(FormatDateTime, "2022-04-01 01:00:00")
	t2, _ := time.Parse(FormatDateTime, "2022-04-01 02:00:00")
	data := url.Values{
		"title":           {"title"},
		"started_at":      {t1.Format(FormatDateTime)},
		"finished_at":     {t2.Format(FormatDateTime)},
		"description":     {"description"},
		"user_id":         {"1"},
		"notify_interval": {"1h"},
	}
	createUrl := "/events/create/"
	r := s.performPost(createUrl, data)
	require.True(s.T(), r.Success)
	defer func() {
		deleteValues := url.Values{
			"id": {r.Data.Event.ID},
		}
		deleteUrl := "/events/delete/?" + deleteValues.Encode()
		r4 := s.performGet(deleteUrl)
		require.True(s.T(), r4.Success)
	}()
	<-time.After(2 * getSelectInterval(s))
	{
		getValues := url.Values{
			"id": {r.Data.Event.ID},
		}
		getUrl := "/events/get/?" + getValues.Encode()
		r1 := s.performGet(getUrl)
		require.True(s.T(), r1.Data.Event.Sent)
	}
}

func (s *SchedulerTestSuite) decodeResponse(resp *http.Response) response {
	var res response
	err := json.NewDecoder(resp.Body).Decode(&res)
	require.Nil(s.T(), err)
	return res
}

func (s *SchedulerTestSuite) performPost(uri string, values url.Values) response {
	fullUrl := s.url + uri
	resp, err := http.PostForm(fullUrl, values)
	require.Nil(s.T(), err)
	return s.decodeResponse(resp)
}

func (s *SchedulerTestSuite) performGet(uri string) response {
	fullUrl := s.url + uri
	resp, err := http.Get(fullUrl)
	require.Nil(s.T(), err)
	return s.decodeResponse(resp)
}

func (s *SchedulerTestSuite) TearDownSuite() {
}

func TestScheduler(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}

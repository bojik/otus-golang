//go:build integration
// +build integration

package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HTTPTestSuite struct {
	suite.Suite
	url string
}

func (s *HTTPTestSuite) SetupSuite() {
	s.url = getHTTPUrl(s)
}

func (s *HTTPTestSuite) TestCRUD() {
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
	require.Equal(s.T(), "title", r.Data.Event.Title)

	{
		// trying to send request again. expected error time is busy
		r := s.performPost(createUrl, data)
		require.False(s.T(), r.Success)
		e := "event date is busy"
		require.Equal(s.T(), e, r.Error[0:len(e)])
	}

	getValues := url.Values{
		"id": {r.Data.Event.ID},
	}
	getUrl := "/events/get/?" + getValues.Encode()
	r1 := s.performGet(getUrl)
	require.True(s.T(), r1.Success)
	require.Equal(s.T(), r.Data.Event.ID, r1.Data.Event.ID)
	require.Equal(s.T(), data["title"][0], r.Data.Event.Title)

	editUrl := "/events/edit/"
	editValues := url.Values{
		"id":    {r.Data.Event.ID},
		"title": {"title1"},
	}
	r2 := s.performPost(editUrl, editValues)
	require.True(s.T(), r2.Success)
	r3 := s.performGet(getUrl)
	require.True(s.T(), r3.Success)
	require.Equal(s.T(), "title1", r3.Data.Event.Title)

	deleteValues := url.Values{
		"id": {r.Data.Event.ID},
	}
	deleteUrl := "/events/delete/?" + deleteValues.Encode()
	r4 := s.performGet(deleteUrl)
	require.True(s.T(), r4.Success)
}

func (s *HTTPTestSuite) TestList() {
	st, _ := time.Parse(FormatDateTime, "2022-05-01 01:00:00")
	ft, _ := time.Parse(FormatDateTime, "2022-05-01 02:00:00")
	const (
		Days0  = 0
		Days3  = 3 * 24 * time.Hour
		Days6  = 6 * 24 * time.Hour
		Days15 = 15 * 24 * time.Hour
		Days25 = 25 * 24 * time.Hour
		Days30 = 30 * 24 * time.Hour
		Days31 = 31 * 24 * time.Hour
	)
	durs := []time.Duration{Days0, Days3, Days6, Days15, Days25, Days30, Days31}
	ids := map[time.Duration]string{}
	for _, d := range durs {
		values := url.Values{
			"title":           {fmt.Sprintf("%d", d)},
			"started_at":      {st.Add(d).Format(FormatDateTime)},
			"finished_at":     {ft.Add(d).Format(FormatDateTime)},
			"description":     {"description"},
			"user_id":         {"1"},
			"notify_interval": {"1h"},
		}
		createUrl := "/events/create/"
		r := s.performPost(createUrl, values)
		require.True(s.T(), r.Success)
		ids[d] = r.Data.Event.ID
	}
	defer func() {
		// delete all
		for _, id := range ids {
			deleteValues := url.Values{
				"id": {id},
			}
			deleteUrl := "/events/delete/?" + deleteValues.Encode()
			r4 := s.performGet(deleteUrl)
			require.True(s.T(), r4.Success)
		}
	}()
	{
		dayList := s.performListRequest("/events/list/day/", st)
		require.True(s.T(), dayList.Success)
		require.Len(s.T(), dayList.Data.Events, 1)
		require.Equal(s.T(), dayList.Data.Events[0].ID, ids[0])
		require.Equal(s.T(), dayList.Data.Events[0].Title, "0")
	}
	{
		weekList := s.performListRequest("/events/list/week/", st)
		require.True(s.T(), weekList.Success)
		require.Len(s.T(), weekList.Data.Events, 3)
		weekDurations := []time.Duration{Days0, Days3, Days6}
		for i, d := range weekDurations {
			require.Equal(s.T(), fmt.Sprintf("%d", d), weekList.Data.Events[i].Title)
			require.Equal(s.T(), ids[d], weekList.Data.Events[i].ID)
		}
	}
	{
		monthList := s.performListRequest("/events/list/month/", st)
		require.True(s.T(), monthList.Success)
		require.Len(s.T(), monthList.Data.Events, 6)
		weekDurations := []time.Duration{Days0, Days3, Days6, Days15, Days25, Days30}
		for i, d := range weekDurations {
			require.Equal(s.T(), fmt.Sprintf("%d", d), monthList.Data.Events[i].Title)
			require.Equal(s.T(), ids[d], monthList.Data.Events[i].ID)
		}
	}
}

func (s *HTTPTestSuite) performListRequest(u string, t time.Time) listResponse {
	listValues := url.Values{
		"date": {t.Format(FormatDate)},
	}
	listUrl := s.url + u + "?" + listValues.Encode()
	resp, err := http.Get(listUrl)
	require.Nil(s.T(), err)
	var res listResponse
	err1 := json.NewDecoder(resp.Body).Decode(&res)
	require.Nil(s.T(), err1)
	return res
}

func (s *HTTPTestSuite) decodeResponse(resp *http.Response) response {
	var res response
	err := json.NewDecoder(resp.Body).Decode(&res)
	require.Nil(s.T(), err)
	return res
}

func (s *HTTPTestSuite) performPost(uri string, values url.Values) response {
	fullUrl := s.url + uri
	resp, err := http.PostForm(fullUrl, values)
	require.Nil(s.T(), err)
	return s.decodeResponse(resp)
}

func (s *HTTPTestSuite) performGet(uri string) response {
	fullUrl := s.url + uri
	resp, err := http.Get(fullUrl)
	require.Nil(s.T(), err)
	return s.decodeResponse(resp)
}

func (s *HTTPTestSuite) TearDownSuite() {
}

func TestHttp(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}

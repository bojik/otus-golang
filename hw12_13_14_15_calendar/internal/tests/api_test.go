//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/pkg/calendarpb"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ApiTestSuite struct {
	suite.Suite
	url string
}

func (s *ApiTestSuite) SetupSuite() {
	s.url = getApiUrl(s)
}

func (s *ApiTestSuite) TestApi() {
	conn, err := grpc.Dial(s.url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(s.T(), err)
	defer func() {
		_ = conn.Close()
	}()
	t1, _ := time.Parse(FormatDateTime, "2022-12-01 10:00:00")
	t2, _ := time.Parse(FormatDateTime, "2022-12-01 11:00:00")
	client := calendarpb.NewCalendarClient(conn)
	e, err := client.InsertEvent(context.Background(), &calendarpb.Event{
		Title:          "title",
		StartedAt:      timestamppb.New(t1),
		FinishedAt:     timestamppb.New(t2),
		Description:    "description",
		UserId:         1,
		NotifyInterval: durationpb.New(time.Hour),
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), "title", e.Title)
	require.NotEmpty(s.T(), e.Id)
	defer func() {
		_, err = client.DeleteEvent(context.Background(), &calendarpb.Id{Id: e.Id})
		require.Nil(s.T(), err)
	}()
}

func (s *ApiTestSuite) TearDownSuite() {
}

func TestApi(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

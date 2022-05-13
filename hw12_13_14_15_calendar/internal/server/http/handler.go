package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"golang.org/x/xerrors"
)

type JSONResponse struct {
	Success bool
	Error   string
	Data    map[string]interface{}
}

const (
	EventFieldTitle          = "title"
	EventFieldStartedAt      = "started_at"
	EventFieldFinishedAt     = "finished_at"
	EventFieldDescription    = "description"
	EventFieldUserID         = "user_id"
	EventFieldNotifyInterval = "notify_interval"
)

// IndexHandler http handler for index page.
type IndexHandler struct {
	logg logger.Logger
}

func (d *IndexHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("<h1>It works</h1>"))
	if err != nil {
		d.logg.Error(err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

var _ http.Handler = (*IndexHandler)(nil)

func NewIndexHandler(logg logger.Logger) *IndexHandler {
	return &IndexHandler{
		logg: logg,
	}
}

// CreateEventHandler http handler for user creation.
type CreateEventHandler struct {
	a *app.App
	l logger.Logger
}

func (c *CreateEventHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	e, err := fillNewEventFromRequest(request)
	if err != nil {
		writeError(err, writer)
		return
	}
	id, err := c.a.CreateEvent(request.Context(), *e)
	if err != nil {
		writeError(err, writer)
		return
	}
	event, err := c.a.FindByID(request.Context(), id)
	if err != nil {
		writeError(err, writer)
		return
	}
	writeSuccess(writer, map[string]interface{}{"Event": event})
}

var _ http.Handler = (*CreateEventHandler)(nil)

func NewCreateEventHandler(a *app.App, l logger.Logger) *CreateEventHandler {
	return &CreateEventHandler{
		a: a,
		l: l,
	}
}

// UpdateEventHandler http handler for user updating.
type UpdateEventHandler struct {
	a *app.App
	l logger.Logger
}

func (c *UpdateEventHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	data, err := parseRequest(request)
	if err != nil {
		writeError(err, writer)
		return
	}
	ids, ok := data["id"]
	if !ok {
		writeError(fmt.Errorf("%w: %s", ErrFieldIsNotDefined, "id"), writer)
		return
	}
	event, err := c.a.FindByID(request.Context(), ids[0])
	if err != nil {
		writeError(err, writer)
		return
	}
	if err := fillEventForEditFromRequest(request, event); err != nil {
		writeError(err, writer)
		return
	}
	if err := c.a.UpdateEvent(request.Context(), *event); err != nil {
		writeError(err, writer)
		return
	}
	writeSuccess(writer, nil)
}

var _ http.Handler = (*UpdateEventHandler)(nil)

func NewUpdateEventHandler(a *app.App, l logger.Logger) *UpdateEventHandler {
	return &UpdateEventHandler{
		a: a,
		l: l,
	}
}

// GetEventHandler http handler for obtaining user information.
type GetEventHandler struct {
	a *app.App
	l logger.Logger
}

func (c GetEventHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	data, err := parseRequest(request)
	if err != nil {
		writeError(err, writer)
		return
	}
	ids, ok := data["id"]
	if !ok {
		writeError(fmt.Errorf("%w: %s", ErrFieldIsNotDefined, "id"), writer)
		return
	}
	event, err := c.a.FindByID(request.Context(), ids[0])
	if err != nil {
		writeError(err, writer)
		return
	}
	writeSuccess(writer, map[string]interface{}{"Event": event})
}

var _ http.Handler = (*GetEventHandler)(nil)

func NewGetEventHandler(a *app.App, l logger.Logger) *GetEventHandler {
	return &GetEventHandler{
		a: a,
		l: l,
	}
}

// DeleteEventHandler http handler for user creation.
type DeleteEventHandler struct {
	a *app.App
	l logger.Logger
}

func (c DeleteEventHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	data, err := parseRequest(request)
	if err != nil {
		writeError(err, writer)
		return
	}
	ids, ok := data["id"]
	if !ok {
		writeError(fmt.Errorf("%w: %s", ErrFieldIsNotDefined, "id"), writer)
		return
	}
	event, err := c.a.DeleteByID(request.Context(), ids[0])
	if err != nil {
		writeError(err, writer)
		return
	}
	writeSuccess(writer, map[string]interface{}{"Event": event})
}

var _ http.Handler = (*DeleteEventHandler)(nil)

func NewDeleteEventHandler(a *app.App, l logger.Logger) *DeleteEventHandler {
	return &DeleteEventHandler{
		a: a,
		l: l,
	}
}

// FindEventsHandler http handler for user selection.
type FindEventsHandler struct {
	a *app.App
	l logger.Logger
}

func (c FindEventsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	data, err := parseRequest(request)
	if err != nil {
		writeError(err, writer)
		return
	}
	from, err := parseTime(data, "from_date")
	if err != nil {
		writeError(err, writer)
		return
	}
	to, err := parseTime(data, "to_date")
	if err != nil {
		writeError(err, writer)
		return
	}
	events, err := c.a.FindEventsByInterval(request.Context(), *from, *to)
	if err != nil {
		writeError(err, writer)
		return
	}
	writeSuccess(writer, map[string]interface{}{"Events": events})
}

var _ http.Handler = (*FindEventsHandler)(nil)

func NewFindEventsHandler(a *app.App, l logger.Logger) *FindEventsHandler {
	return &FindEventsHandler{
		a: a,
		l: l,
	}
}

// FindEventsDayHandler http handler for user selection for day.
type FindEventsDayHandler struct {
	a *app.App
	l logger.Logger
}

func (c FindEventsDayHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	serveDateHTTP(writer, request, c.a.FindDayEvents)
}

var _ http.Handler = (*FindEventsDayHandler)(nil)

func NewFindEventsDayHandler(a *app.App, l logger.Logger) *FindEventsDayHandler {
	return &FindEventsDayHandler{
		a: a,
		l: l,
	}
}

// FindEventsWeekHandler http handler for user selection for Week.
type FindEventsWeekHandler struct {
	a *app.App
	l logger.Logger
}

func (c FindEventsWeekHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	serveDateHTTP(writer, request, c.a.FindWeekEvents)
}

var _ http.Handler = (*FindEventsWeekHandler)(nil)

func NewFindEventsWeekHandler(a *app.App, l logger.Logger) *FindEventsWeekHandler {
	return &FindEventsWeekHandler{
		a: a,
		l: l,
	}
}

// FindEventsMonthHandler http handler for user selection for Month.
type FindEventsMonthHandler struct {
	a *app.App
	l logger.Logger
}

func (c FindEventsMonthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	serveDateHTTP(writer, request, c.a.FindMonthEvents)
}

var _ http.Handler = (*FindEventsMonthHandler)(nil)

func NewFindEventsMonthHandler(a *app.App, l logger.Logger) *FindEventsMonthHandler {
	return &FindEventsMonthHandler{
		a: a,
		l: l,
	}
}

type findFuncType func(context.Context, time.Time) ([]*app.Event, error)

func serveDateHTTP(writer http.ResponseWriter, request *http.Request, findFunc findFuncType) {
	data, err := parseRequest(request)
	if err != nil {
		writeError(err, writer)
		return
	}
	date, err := parseTime(data, "date")
	if err != nil {
		writeError(err, writer)
		return
	}
	events, err := findFunc(request.Context(), *date)
	if err != nil {
		writeError(err, writer)
		return
	}
	writeSuccess(writer, map[string]interface{}{"Events": events})
}

func writeError(err error, writer http.ResponseWriter) {
	resp := JSONResponse{
		Success: false,
		Error:   err.Error(),
	}
	bytes, _ := json.Marshal(resp)
	writer.Header().Set("content-type", "application-json")
	writer.Write(bytes)
}

func writeSuccess(writer http.ResponseWriter, data map[string]interface{}) {
	resp := JSONResponse{
		Success: true,
		Data:    data,
	}
	bytes, _ := json.Marshal(resp)
	writer.Header().Set("content-type", "application-json")
	_, _ = writer.Write(bytes)
}

func fillNewEventFromRequest(r *http.Request) (*app.Event, error) {
	data, err := parseRequest(r)
	if err != nil {
		return nil, err
	}
	e := &app.Event{}
	{
		strs, ok := data[EventFieldTitle]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrFieldIsNotDefined, EventFieldTitle)
		}
		e.Title = strs[0]
	}
	{
		userID, err := parseInt(data, EventFieldUserID)
		if err != nil {
			return nil, err
		}
		e.UserID = userID
	}
	{
		t, err := parseTime(data, EventFieldStartedAt)
		if err != nil {
			return nil, err
		}
		e.StartedAt = *t
	}
	{
		t, err := parseTime(data, EventFieldFinishedAt)
		if err != nil {
			return nil, err
		}
		e.FinishedAt = *t
	}
	{
		strs, ok := data[EventFieldDescription]
		if ok {
			e.Description = strs[0]
		}
	}
	{
		d, err := parseDuration(data, EventFieldNotifyInterval)
		if err != nil && !xerrors.Is(err, ErrFieldIsNotDefined) {
			return nil, err
		}
		if !xerrors.Is(err, ErrFieldIsNotDefined) {
			e.NotifyInterval = *d
		}
	}
	return e, nil
}

func fillEventForEditFromRequest(r *http.Request, e *app.Event) error {
	data, err := parseRequest(r)
	if err != nil {
		return err
	}
	{
		strs, ok := data[EventFieldTitle]
		if ok {
			e.Title = strs[0]
		}
	}
	{
		t, err := parseTime(data, EventFieldStartedAt)
		if err != nil && !xerrors.Is(err, ErrFieldIsNotDefined) {
			return err
		}
		if !xerrors.Is(err, ErrFieldIsNotDefined) {
			e.StartedAt = *t
		}
	}
	{
		t, err := parseTime(data, EventFieldFinishedAt)
		if err != nil && !xerrors.Is(err, ErrFieldIsNotDefined) {
			return err
		}
		if !xerrors.Is(err, ErrFieldIsNotDefined) {
			e.FinishedAt = *t
		}
	}
	{
		strs, ok := data[EventFieldDescription]
		if ok {
			e.Description = strs[0]
		}
	}
	{
		userID, err := parseInt(data, EventFieldUserID)
		if err != nil && !xerrors.Is(err, ErrFieldIsNotDefined) {
			return err
		}
		if !xerrors.Is(err, ErrFieldIsNotDefined) {
			e.UserID = userID
		}
	}
	{
		d, err := parseDuration(data, EventFieldNotifyInterval)
		if err != nil && !xerrors.Is(err, ErrFieldIsNotDefined) {
			return err
		}
		if !xerrors.Is(err, ErrFieldIsNotDefined) {
			e.NotifyInterval = *d
		}
	}
	return nil
}

func parseRequest(r *http.Request) (url.Values, error) {
	data := r.URL.Query()
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		data = r.PostForm
	}
	return data, nil
}

func parseTime(data url.Values, key string) (*time.Time, error) {
	strs, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrFieldIsNotDefined, key)
	}
	t, err := time.Parse(time.RFC3339, strs[0])
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseDuration(data url.Values, key string) (*time.Duration, error) {
	strs, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrFieldIsNotDefined, key)
	}
	d, err := time.ParseDuration(strs[0])
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func parseInt(data url.Values, key string) (int, error) {
	strs, ok := data[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrFieldIsNotDefined, key)
	}
	n, err := strconv.Atoi(strs[0])
	if err != nil {
		return 0, err
	}
	return n, nil
}

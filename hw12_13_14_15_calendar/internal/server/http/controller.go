package internalhttp

import (
	"net/http"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
)

type DefaultController struct {
	logg logger.Logger
}

func (d *DefaultController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	time.Sleep(time.Second + 500*time.Millisecond)
	_, err := writer.Write([]byte("<h1>It works</h1>"))
	if err != nil {
		d.logg.Error(err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

var _ http.Handler = (*DefaultController)(nil)

func NewDefaultController(logg logger.Logger) *DefaultController {
	return &DefaultController{
		logg: logg,
	}
}

package logger

import (
	"errors"
	"fmt"
)

type Option interface {
	getOption() interface{}
}

type OptionMinLevel struct {
	minLevel Level
}

func (o *OptionMinLevel) getOption() interface{} {
	return o.minLevel
}

var _ Option = (*OptionMinLevel)(nil)

func NewOptionMinLevel(lvl string) (*OptionMinLevel, error) {
	lvls := []Level{DEBUG, INFO, ERROR}
	for _, l := range lvls {
		if l.String() == lvl {
			return &OptionMinLevel{minLevel: l}, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Unknown log Level: %s", lvl))
}

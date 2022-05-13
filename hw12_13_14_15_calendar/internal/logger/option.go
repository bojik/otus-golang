package logger

import (
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
	return nil, fmt.Errorf("unknown log Level: %s", lvl)
}

package main

import (
	"errors"

	"golang.org/x/xerrors"
)

var ErrIsNotConnected = errors.New("client is not connected, please, call Connect()")

type TelnetError struct {
	err  error
	addr string
}

func (t TelnetError) Unwrap() error {
	return t.err
}

func (t TelnetError) Error() string {
	return t.err.Error()
}

var (
	_ error           = (*TelnetError)(nil)
	_ xerrors.Wrapper = (*TelnetError)(nil)
)

func newTelnetError(addr string, err error) error {
	return TelnetError{
		addr: addr,
		err:  err,
	}
}

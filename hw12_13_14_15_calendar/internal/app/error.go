package app

import "golang.org/x/xerrors"

var (
	ErrDateBusy = xerrors.New("event date is busy")
)

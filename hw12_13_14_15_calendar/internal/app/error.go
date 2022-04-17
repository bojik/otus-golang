package app

import "golang.org/x/xerrors"

var (
	ErrDateBusy      = xerrors.New("event date is busy")
	ErrRequiredField = xerrors.New("is required field")
	ErrInvalidDate   = xerrors.New("invalid date")
)

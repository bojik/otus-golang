package main

import (
	"errors"

	"golang.org/x/xerrors"
)

type CopyError struct {
	Path string
	Err  error
}

func (e *CopyError) Error() string {
	return e.Path + " - " + e.Err.Error()
}

func (e *CopyError) Unwrap() error {
	return e.Err
}

var (
	_ error           = (*CopyError)(nil)
	_ xerrors.Wrapper = (*CopyError)(nil)
)

func NewCopyError(path string, err error) *CopyError {
	return &CopyError{
		Path: path,
		Err:  err,
	}
}

var (
	ErrFileNotExists         = errors.New("file does not exist")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEmptyDestinationPath  = errors.New("destination path is empty")
	ErrEmptySourcePath       = errors.New("source path is empty")
	ErrInvalidOffset         = errors.New("invalid offset")
	ErrInvalidLimit          = errors.New("invalid limit")
	ErrIsDirectory           = errors.New("source path is a directory")
)

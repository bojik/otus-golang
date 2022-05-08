package memorystorage

import "errors"

var (
	ErrEventAlreadyInserted = errors.New("event is already inserted")
	ErrEventNotFound        = errors.New("event not found")
)

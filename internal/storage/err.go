package storage

import (
	"errors"
)

var (
	ErrNotFound = errors.New("task not found")
	ErrNoUpdate = errors.New("nothing to update")
)

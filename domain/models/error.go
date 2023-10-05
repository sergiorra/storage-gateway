package models

import "fmt"

type (
	ErrNotFound struct {
		value string
	}

	ErrNotValid struct {
		value string
	}

	ErrNotAvailable struct {
		value string
	}
)

const (
	object        = "object"
	objectStorage = "object storage"
	objectID      = "object ID"
)

var (
	ErrObjectNotFound            = NewErrNotFound(object)
	ErrObjectIDNotValid          = NewErrObjectIDNotValid(objectID)
	ErrObjectStorageNotAvailable = NewErrObjectStorageNotAvailable(objectStorage)
)

func NewErrNotFound(value string) *ErrNotFound {
	return &ErrNotFound{value}
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found", err.value)
}

func NewErrObjectIDNotValid(value string) *ErrNotValid {
	return &ErrNotValid{value}
}

func (err ErrNotValid) Error() string {
	return fmt.Sprintf("%s not valid", err.value)
}

func NewErrObjectStorageNotAvailable(value string) *ErrNotAvailable {
	return &ErrNotAvailable{value}
}

func (err ErrNotAvailable) Error() string {
	return fmt.Sprintf("%s not available", err.value)
}

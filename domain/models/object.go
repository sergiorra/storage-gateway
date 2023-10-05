package models

import (
	"io"
)

type Object struct {
	ID          ObjectID
	Content     io.Reader
	ContentType string
	Size        int64
}

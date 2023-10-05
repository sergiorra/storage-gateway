package models

import "regexp"

type ObjectID string

func (id ObjectID) IsValidID() bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]{1,32}$`).MatchString(id.Value())
}

func (id ObjectID) Value() string {
	return string(id)
}

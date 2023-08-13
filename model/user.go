package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Name     string
	Password string
}

package domain

import (
	"github.com/google/uuid"
)

type Status string

const (
	StatusActive Status = "Active"
	StatusFrozen Status = "Fronzen"
	StatusClosed Status = "Closed"
)

type Account struct {
	ID      uuid.UUID
	No      uint64
	Name    string
	Balance uint64
	Status  Status
}

package domain

import (
	"github.com/google/uuid"
	"time"
)

type EventType string

const (
	EventAccountCreated EventType = "AccountCreated"
	EventMoneyDeposited EventType = "MoneyDeposited"
	EventMoneyWithdrawn EventType = "MoneyWithdrawn"
)

type Event struct {
	ID        uuid.UUID
	AccountId uuid.UUID
	Type      EventType
	Timestamp time.Time
}

type AccountCreated struct {
	Event
	Name   string
	No     uint64
	Amount uint64
}

type MoneyDeposited struct {
	Event
	Amount uint64
}

type MoneyWithdrawn struct {
	Event
	Amount uint64
}

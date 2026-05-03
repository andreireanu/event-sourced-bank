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
	AccountID uuid.UUID
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

type DomainEvent interface {
	isEvent()
	GetAccountID() uuid.UUID
}

func (event Event) isEvent() {}

func (event Event) GetAccountID() uuid.UUID {
	return event.AccountID
}

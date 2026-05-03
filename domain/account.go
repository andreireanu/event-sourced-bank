package domain

import (
	"errors"
	"github.com/google/uuid"
)

type Status string

const (
	StatusActive Status = "Active"
	StatusFrozen Status = "Frozen"
	StatusClosed Status = "Closed"
)

type Account struct {
	ID      uuid.UUID
	No      uint64
	Name    string
	Balance uint64
	Status  Status
}

func (account *Account) applyEvent(event DomainEvent) error {
	switch e := event.(type) {
	case AccountCreated:
		account.ID = e.ID
		account.No = e.No
		account.Name = e.Name
		account.Balance = e.Amount
		account.Status = StatusActive
	case MoneyDeposited:
		account.Balance += e.Amount
	case MoneyWithdrawn:
		if e.Amount > account.Balance {
			return errors.New("Insufficient balance")
		}
		account.Balance -= e.Amount
	}
	return nil
}

func LoadAccount(events []DomainEvent) (Account, error) {
	var account Account
	for _, ev := range events {
		result := account.applyEvent(ev)
		if result != nil {
			return Account{}, result
		}
	}
	return account, nil
}

package main

import (
	"bank/domain"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func main() {
	events := make([]domain.DomainEvent, 3)

	accountCreated := domain.AccountCreated{
		Event: domain.Event{
			ID:        uuid.New(),
			AccountID: uuid.New(),
			Type:      domain.EventAccountCreated,
			Timestamp: time.Now(),
		},
		Name:   "Alice",
		No:     1,
		Amount: 0,
	}

	moneyDeposited := domain.MoneyDeposited{
		Event: domain.Event{
			ID:        uuid.New(),
			AccountID: accountCreated.AccountID,
			Type:      domain.EventMoneyDeposited,
			Timestamp: time.Now(),
		},
		Amount: 100,
	}

	moneyWithdrawn := domain.MoneyWithdrawn{
		Event: domain.Event{
			ID:        uuid.New(),
			AccountID: accountCreated.AccountID,
			Type:      domain.EventMoneyWithdrawn,
			Timestamp: time.Now(),
		},
		Amount: 50,
	}

	events = append(events, accountCreated)
	events = append(events, moneyDeposited)
	events = append(events, moneyWithdrawn)

	account, err := domain.LoadAccount(events)
	if err == nil {
		fmt.Println("&v+\n", account)
	}
}

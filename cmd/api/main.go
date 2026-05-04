package main

import (
	"bank/domain"
	"bank/store"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func main() {
	memoryStore := store.NewMemoryStore()
	aliceID := uuid.New()

	accountCreated := domain.AccountCreated{
		Event: domain.Event{
			ID:        uuid.New(),
			AccountID: aliceID,
			Type:      domain.EventAccountCreated,
			Timestamp: time.Now(),
		},
		Name:   "Alice",
		No:     1,
		Amount: 0,
	}
	_ = memoryStore.Save(accountCreated)

	moneyDeposited := domain.MoneyDeposited{
		Event: domain.Event{
			ID:        uuid.New(),
			AccountID: aliceID,
			Type:      domain.EventMoneyDeposited,
			Timestamp: time.Now(),
		},
		Amount: 100,
	}
	_ = memoryStore.Save(moneyDeposited)

	moneyWithdrawn := domain.MoneyWithdrawn{
		Event: domain.Event{
			ID:        uuid.New(),
			AccountID: aliceID,
			Type:      domain.EventMoneyWithdrawn,
			Timestamp: time.Now(),
		},
		Amount: 50,
	}
	_ = memoryStore.Save(moneyWithdrawn)

	events, ok := memoryStore.Load(aliceID)
	if ok != nil {
		fmt.Println("ID not found")
	}

	account, err := domain.LoadAccount(events)
	if err == nil {
		fmt.Println(account)
	} else {
		fmt.Println(err)
	}

	bobId := uuid.New()
	_, ok = memoryStore.Load(bobId)
	if ok != nil {
		fmt.Println("ID not found")
	}

}

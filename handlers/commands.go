package handlers

import (
	"bank/domain"
	"bank/store"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type CreateAccountCommand struct {
	Name   string `json:"name"`
	Amount uint64 `json:"amount"`
}

type DepositCommand struct {
	AccountID uuid.UUID
	Amount    uint64
}

type WithdrawCommand struct {
	AccountID uuid.UUID
	Amount    uint64
}
type CommandHandler struct {
	MemoryStore *store.MemoryStore
}

func NewCommandHandler(memoryStore *store.MemoryStore) CommandHandler {
	return CommandHandler{MemoryStore: memoryStore}
}

func (comHandler *CommandHandler) CreateAccount(cmd CreateAccountCommand) (uuid.UUID, uint64, error) {
	accountID := uuid.New()
	accountNo := rand.Uint64() % 100_000_000
	accountCreated := domain.AccountCreated{
		Event: domain.Event{
			ID:        uuid.New(),
			AccountID: accountID,
			Type:      domain.EventAccountCreated,
			Timestamp: time.Now(),
		},
		Name:   cmd.Name,
		No:     accountNo,
		Amount: cmd.Amount,
	}
	ok := comHandler.MemoryStore.Save(accountCreated)
	return accountID, accountNo, ok
}

func (comHandler *CommandHandler) Deposit(cmd DepositCommand) error {
	events, ok := comHandler.MemoryStore.Load(cmd.AccountID)
	if ok != nil {
		return ok
	}

	account, ok := domain.LoadAccount(events)
	if ok != nil {
		return ok
	}

	ok = account.Deposit(cmd.Amount)
	if ok != nil {
		return ok
	}

	moneyDeposited := domain.MoneyDeposited{
		Event: domain.Event{
			ID:        uuid.New(),
			AccountID: cmd.AccountID,
			Type:      domain.EventMoneyDeposited,
			Timestamp: time.Now(),
		},
		Amount: cmd.Amount,
	}
	ok = comHandler.MemoryStore.Save(moneyDeposited)
	return ok
}

func (comHandler *CommandHandler) Withdraw(cmd WithdrawCommand) error {
	events, ok := comHandler.MemoryStore.Load(cmd.AccountID)
	if ok != nil {
		return ok
	}

	account, ok := domain.LoadAccount(events)
	if ok != nil {
		return ok
	}

	ok = account.Withdraw(cmd.Amount)
	if ok != nil {
		return ok
	}

	moneyWithdrawn := domain.MoneyWithdrawn{
		Event: domain.Event{
			ID:        uuid.New(),
			AccountID: cmd.AccountID,
			Type:      domain.EventMoneyWithdrawn,
			Timestamp: time.Now(),
		},
		Amount: cmd.Amount,
	}
	ok = comHandler.MemoryStore.Save(moneyWithdrawn)
	return ok
}

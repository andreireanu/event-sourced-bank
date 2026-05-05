package handlers

import (
	"bank/domain"
	"bank/store"
	"github.com/google/uuid"
)

type AccountQuery struct {
	AccountID uuid.UUID
}

type QueryHandler struct {
	MemoryStore *store.MemoryStore
}

func NewQueryHandler(memoryStore *store.MemoryStore) QueryHandler {
	return QueryHandler{MemoryStore: memoryStore}
}
func (queryHandler *QueryHandler) GetAccount(accountQuery AccountQuery) (domain.Account, error) {
	events, ok := queryHandler.MemoryStore.Load(accountQuery.AccountID)
	if ok != nil {
		return domain.Account{}, ok
	}
	account, ok := domain.LoadAccount(events)
	if ok != nil {
		return domain.Account{}, ok
	}
	return account, nil
}

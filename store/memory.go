package store

import (
	"bank/domain"
	"github.com/google/uuid"
)

type MemoryStore struct {
	Map map[uuid.UUID][]domain.DomainEvent
}

func NewMemoryStore() MemoryStore {
	newMap := MemoryStore{Map: make(map[uuid.UUID][]domain.DomainEvent, 0)}
	return newMap
}

func (memoryStore *MemoryStore) Save(event domain.DomainEvent) error {
	id := event.GetAccountID()
	mp, ok := memoryStore.Map[id]
	if ok {
		memoryStore.Map[id] = append(mp, event)
	} else {
		memoryStore.Map[id] = []domain.DomainEvent{event}
	}
	return nil
}

package main

import (
	// "bank/domain"
	"bank/handlers"
	"bank/store"
	// "fmt"
	// "github.com/google/uuid"
	// "time"
)

func main() {

	memoryStore := store.NewMemoryStore()
	commandHandler := handlers.NewCommandHandler(&memoryStore)
	queryHandler := handlers.NewQueryHandler(&memoryStore)

	server := NewServer(&commandHandler, &queryHandler)
	server.Start()
}

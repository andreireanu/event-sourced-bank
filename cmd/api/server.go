package main

import (
	"bank/handlers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type Server struct {
	commands *handlers.CommandHandler
	queries  *handlers.QueryHandler
}

func NewServer(commands *handlers.CommandHandler, queries *handlers.QueryHandler) *Server {
	return &Server{
		commands,
		queries,
	}
}

func (server *Server) CreateAccount(w http.ResponseWriter, req *http.Request) {
	requestID := uuid.NewString()
	log.Printf("[%s] Creating account", requestID)

	method := req.Method
	if method != "POST" {
		log.Printf("[%s] Mehtod not allowed", requestID)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body_decoder := json.NewDecoder(req.Body)
	var cmd handlers.CreateAccountCommand
	ok := body_decoder.Decode(&cmd)
	if ok != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id, no, ok := server.commands.CreateAccount(cmd)
	if ok != nil {
		log.Printf("[%s] Cannot create account: %s", requestID, ok)
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Received message",
		"body":    "Created account with id " + id.String() + "and number " + strconv.FormatUint(no, 10),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
	log.Printf("[%s] Account created", requestID)
}
func (server *Server) Start() {
	http.HandleFunc("/accounts", server.CreateAccount)
	fmt.Println("Server starting on :8090")
	err := http.ListenAndServe(":8090", nil)
	fmt.Println("Server stopped:", err)
}

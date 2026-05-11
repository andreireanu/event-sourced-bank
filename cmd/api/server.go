package main

import (
	"bank/domain"
	"bank/handlers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
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
		log.Printf("[%s] Mehtod not allowed: %s", requestID, method)
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
		"body":    "Created account with id " + id.String() + " and number " + strconv.FormatUint(no, 10),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
	log.Printf("[%s] Account created with id %s", requestID, id)
}

func (server *Server) Deposit(w http.ResponseWriter, req *http.Request) {
	requestID := uuid.NewString()
	id_string := req.PathValue("id")
	log.Printf("[%s] Attempting to deposit for account id %s", requestID, id_string)

	method := req.Method
	if method != "POST" {
		log.Printf("[%s] Mehtod not allowed: %s", requestID, method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body_decoder := json.NewDecoder(req.Body)
	var cmd handlers.DepositCommand
	ok := body_decoder.Decode(&cmd)
	if ok != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	cmd.AccountID, ok = uuid.Parse(id_string)
	if ok != nil {
		log.Printf("[%s] Error parsing account id: %s, error: %s", requestID, id_string, ok)
		http.Error(w, "Failed to deposit into account", http.StatusInternalServerError)
		return
	}

	ok = server.commands.Deposit(cmd)
	switch {
	case errors.Is(ok, domain.ErrAmountTooHigh):
		log.Printf("[%s] Error depositing to account id: %s, error: amount too high to withdraw", requestID, id_string)
		http.Error(w, "Failed to deposit to account, amount too high", http.StatusBadRequest)
		return

	case ok != nil:
		log.Printf("[%s] Cannot deposit to account: %s, error: %s", requestID, id_string, ok)
		http.Error(w, "Failed to deposit into account", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Received message",
		"body":    "Deposited to account with id " + id_string + " amount " + strconv.FormatUint(cmd.Amount, 10),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
	log.Printf("[%s] Deposited to account with id %s amount %s", requestID, id_string, strconv.FormatUint(cmd.Amount, 10))
}

func (server *Server) Withdraw(w http.ResponseWriter, req *http.Request) {
	requestID := uuid.New()
	id_string := req.PathValue("id")

	method := req.Method
	if method != "POST" {
		log.Printf("[%s] Method not allowed: %s", requestID, method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cmd handlers.WithdrawCommand

	decoder := json.NewDecoder(req.Body)
	ok := decoder.Decode(&cmd)
	if ok != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id, ok := uuid.Parse(id_string)
	if ok != nil {
		log.Printf("[%s] Error parsing account id: %s, error: %s", requestID, id, ok)
		http.Error(w, "Failed to withdraw from account", http.StatusInternalServerError)
		return
	}
	cmd.AccountID = id

	ok = server.commands.Withdraw(cmd)
	switch {
	case errors.Is(ok, domain.ErrInsufficientFunds):
		log.Printf("[%s] Error withdrawing from account id: %s, error: insufficient funds", requestID, id)
		http.Error(w, "Failed to withdraw from account, insufficient funds", http.StatusBadRequest)
		return
	case ok != nil:
		log.Printf("[%s] Error withdrawing from account id: %s, error: %s", requestID, id, ok)
		http.Error(w, "Failed to withdraw from account", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Received message",
		"body":    "Withdrawn from account with id " + id_string + " amount " + strconv.FormatUint(cmd.Amount, 10),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
	log.Printf("[%s] Withdrawn from account with id %s amount %s", requestID, id, strconv.FormatUint(cmd.Amount, 10))
}

func (server *Server) GetAccount(w http.ResponseWriter, req *http.Request) {
	requestID := uuid.New()

	method := req.Method
	if method != "GET" {
		log.Printf("[%s] Method not allowed: %s", requestID, method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id_string := req.PathValue("id")
	id, ok := uuid.Parse(id_string)
	if ok != nil {
		log.Printf("[%s] Error parsing account id: %s, error: %s", requestID, id, ok)
		http.Error(w, "Failed to withdraw from account", http.StatusInternalServerError)
		return
	}

	query := handlers.AccountQuery{AccountID: id}
	account, ok := server.queries.GetAccount(query)
	if ok != nil {
		log.Printf("[%s] Error loading account: %s, error: %s", requestID, id, ok)
		http.Error(w, "Failed to load account", http.StatusInternalServerError)
		return
	}
	data, ok := json.Marshal(account)
	if ok != nil {
		log.Printf("[%s] Error serializing account: %s, error: %s", requestID, id, ok)
		http.Error(w, "Failed to serialize account", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Received message",
		"body":    "Succesfully loaded account %s" + id_string,
		"data":    string(data),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
	log.Printf("[%s] Queried from account with id %s", requestID, id)
}

func (server *Server) Start() {

	http.HandleFunc("/accounts", server.CreateAccount)
	http.HandleFunc("/accounts/{id}/deposit", server.Deposit)
	http.HandleFunc("/accounts/{id}/withdraw", server.Withdraw)
	http.HandleFunc("/accounts/{id}/query", server.GetAccount)

	fmt.Println("Server starting on :8090")
	err := http.ListenAndServe(":8090", nil)
	fmt.Println("Server stopped:", err)
}

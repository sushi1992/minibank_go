package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sushi1992/minibank_go/internal/account"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type CreateAccountRequest struct {
	ID       string           `json:"id"`
	Owner    string           `json:"owner"`
	Currency account.Currency `json:"currency"`
}

type AmountRequest struct {
	Amount int64 `json:"amount"`
}

type TransferRequest struct {
	FromID string `json:"fromID"`
	ToID   string `json:"toID"`
	Amount int64  `json:"amount"`
}

func withdrawHandler(service *account.Service) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")

		var amountRequest AmountRequest
		err := json.NewDecoder(request.Body).Decode(&amountRequest)
		if err != nil {
			http.Error(writer, "invalid request body", http.StatusBadRequest)
			return
		}

		updatedAccount, err := service.Withdraw(id, amountRequest.Amount)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(writer).Encode(updatedAccount)
		if err != nil {
			return
		}
	}
}

func depositHandler(service *account.Service) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")

		var amountRequest AmountRequest

		err := json.NewDecoder(request.Body).Decode(&amountRequest)
		if err != nil {
			http.Error(writer, "invalid request body", http.StatusBadRequest)
			return
		}

		updatedAccount, err := service.Deposit(id, amountRequest.Amount)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(writer).Encode(updatedAccount)
	}
}

func getAccountHandler(service *account.Service) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")

		retrievedAccount, err := service.GetAccount(id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}

		writer.Header().Set("content-type", "application/json")

		err = json.NewEncoder(writer).Encode(retrievedAccount)
		if err != nil {
			http.Error(writer, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func transferHandler(service *account.Service) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var transferRequest TransferRequest

		err := json.NewDecoder(request.Body).Decode(&transferRequest)
		if err != nil {
			http.Error(writer, "invalid request body", http.StatusBadRequest)
			return
		}

		err = service.Transfer(transferRequest.FromID, transferRequest.ToID, transferRequest.Amount)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}

func createAccountHandler(service *account.Service) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var createRequest CreateAccountRequest

		err := json.NewDecoder(request.Body).Decode(&createRequest)
		if err != nil {
			http.Error(writer, "invalid request body", http.StatusBadRequest)
			return
		}

		createAccount, err := service.CreateAccount(createRequest.ID, createRequest.Owner, createRequest.Currency)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(writer).Encode(createAccount)
		if err != nil {
			return
		}
	}
}

func healthHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	response := HealthResponse{
		Status: "ok",
	}

	json.NewEncoder(writer).Encode(response)
}

func main() {
	repo, err := account.NewCassandraRepository("127.0.0.1")
	if err != nil {
		panic(err)
	}
	service := account.NewService(repo)

	publisher := account.NewKafkaPublisher("localhost:9092")
	service.SetPublished(publisher)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	consumer := account.NewKafkaConsumer("localhost:9092")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := consumer.Consume(ctx)
		if err != nil {
			fmt.Printf("consumer stopped with error: %v\n", err)
		}
	}()

	http.HandleFunc("POST /accounts", createAccountHandler(service))
	http.HandleFunc("GET /health", healthHandler)
	http.HandleFunc("GET /accounts/{id}", getAccountHandler(service))
	http.HandleFunc("POST /accounts/{id}/deposit", depositHandler(service))
	http.HandleFunc("POST /accounts/{id}/withdraw", withdrawHandler(service))
	http.HandleFunc("POST /transfers", transferHandler(service))

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("http server error: %v\n", err)
		}
	}()

	// Wait here until Ctrl+C / SIGTERM
	<-ctx.Done()

	fmt.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		fmt.Printf("http shutdown error: %v\n", err)
	}

	wg.Wait()

	fmt.Print("shutdown complete")
}

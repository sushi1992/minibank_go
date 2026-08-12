package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sushi1992/minibank_go/internal/account"
)

func TestCreateAccountHandler(t *testing.T) {
	repo := account.NewMemoryRepository()
	service := account.NewService(repo)

	handler := createAccountHandler(service)

	body := `{
        "id": "acc-001",
        "owner": "John Doe",
        "currency": "GBP"
    }`

	request := httptest.NewRequest(
		http.MethodPost,
		"/accounts",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	handler(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}

	if recorder.Header().Get("content-type") != "application/json" {
		t.Errorf(
			"expected application/json, got %s",
			recorder.Header().Get("content-type"),
		)
	}

	var createdAccount account.Account

	err := json.NewDecoder(recorder.Body).Decode(&createdAccount)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdAccount.ID != "acc-001" {
		t.Errorf("expected ID acc-001, got %s", createdAccount.ID)
	}

	if createdAccount.Owner != "John Doe" {
		t.Errorf("expected owner John Doe, got %s", createdAccount.Owner)
	}

	if createdAccount.Currency != account.CurrencyGBP {
		t.Errorf("expected GBP, got %s", createdAccount.Currency)
	}
}

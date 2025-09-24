package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/iamwy7/payment-gateway-management/internal/application/service"
	"github.com/iamwy7/payment-gateway-management/internal/domain"
	"github.com/iamwy7/payment-gateway-management/internal/infra/http_api/dtos"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dtos.CreateAccountInput
	// TODO:
	// - Validate Name size and Email regex and etc
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.accountService.CreateAccount(domain.NewAccount(input.Name, input.Email))

	// TODO:
	// - Domain Error Messages
	if err != nil {
		switch err {
		case domain.ErrEmailAlreadyExists:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	output, err := h.accountService.FindByAPIKey(apiKey)
	// TODO:
	// - Domain Error Messages
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

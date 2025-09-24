package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iamwy7/payment-gateway-management/internal/application/service"
	"github.com/iamwy7/payment-gateway-management/internal/domain"
	"github.com/iamwy7/payment-gateway-management/internal/infra/http_api/dtos"
)

type InvoiceHandler struct {
	invoiceService service.InvoiceService
}

func NewInvoiceHandler(invoiceService service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: invoiceService,
	}
}

func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")

	var invoiceInput dtos.CreateInvoiceInput
	err := json.NewDecoder(r.Body).Decode(&invoiceInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // custom error message
		return
	}

	invoiceRequest := dtos.FromInvoiceInput(&invoiceInput)
	invoiceCreated, err := h.invoiceService.Create(apiKey, invoiceRequest)
	// TODO:
	// - Domain Error Messages
	if err != nil {
		switch err {
		case domain.ErrAccountNotFound:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		case domain.ErrUnauthorizedAccess:
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoiceCreated)

}

func (h *InvoiceHandler) GetById(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")

	invoiceId := chi.URLParam(r, "invoiceId")
	if invoiceId == "" {
		http.Error(w, "invoice id is required", http.StatusBadRequest)
	}

	invoice, err := h.invoiceService.GetById(apiKey, invoiceId)
	// TODO:
	// - Domain Error Messages
	if err != nil {
		switch err {
		case domain.ErrInvoiceNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case domain.ErrAccountNotFound:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		case domain.ErrUnauthorizedAccess:
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)

}

func (h *InvoiceHandler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")

	invoices, err := h.invoiceService.ListByAccountApiKey(apiKey)
	// TODO:
	// - Domain Error Messages
	if err != nil {
		switch err {
		case domain.ErrInvoicesNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case domain.ErrAccountNotFound:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoices)
}

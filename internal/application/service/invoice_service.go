package service

import (
	"github.com/iamwy7/payment-gateway-management/internal/domain"
	"github.com/iamwy7/payment-gateway-management/internal/infra/http_api/dtos"
	"github.com/iamwy7/payment-gateway-management/internal/ports"
)

type InvoiceService struct {
	invoiceRepository ports.InvoiceRepository
	accountService    *AccountService
}

func NewInvoiceService(invoiceRepo ports.InvoiceRepository, accountService *AccountService) *InvoiceService {
	return &InvoiceService{
		invoiceRepository: invoiceRepo,
		accountService:    accountService,
	}
}

func (s *InvoiceService) Create(apikey string, invoiceReq *domain.InvoiceRequest) (*domain.Invoice, error) {
	// TODO:
	// - Idempotence Key
	account, err := s.accountService.FindByAPIKey(apikey)
	if err != nil {
		return nil, err
	}

	invoice, err := dtos.ToInvoice(invoiceReq, account.ID)
	if err != nil {
		return nil, err
	}

	if err := invoice.Process(); err != nil {
		return nil, err
	}

	if invoice.Status == domain.StatusApproved {
		account.Balance += invoice.Amount
		err = s.accountService.repository.UpdateBalance(account)
		if err != nil {
			return nil, err
		}
	}

	if err := s.invoiceRepository.Save(invoice); err != nil {
		return nil, err
	}
	return invoice, nil
}

func (s *InvoiceService) GetById(apikey, invoiceId string) (*domain.Invoice, error) {
	invoice, err := s.invoiceRepository.FindByID(invoiceId)
	if err != nil {
		return nil, err
	}

	account, err := s.accountService.repository.FindByAPIKey(apikey)
	if err != nil {
		return nil, err
	}

	if invoice.AccountID != account.ID {
		return nil, domain.ErrUnauthorizedAccess
	}
	return invoice, nil
}

func (s *InvoiceService) ListByAccountId(accountId string) ([]*domain.Invoice, error) {
	invoices, err := s.invoiceRepository.FindByAccountID(accountId)

	if err != nil {
		return nil, err
	}
	return invoices, nil
}

func (s *InvoiceService) ListByAccountApiKey(apikey string) ([]*domain.Invoice, error) {
	account, err := s.accountService.FindByAPIKey(apikey)
	if err != nil {
		return nil, err
	}
	return s.ListByAccountId(account.ID)
}

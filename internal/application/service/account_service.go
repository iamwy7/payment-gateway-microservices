package service

import (
	"github.com/iamwy7/payment-gateway-management/internal/domain"
	"github.com/iamwy7/payment-gateway-management/internal/ports"
)

// AccountService implementa a lógica de negócios para operações com Account
type AccountService struct {
	repository ports.AccountRepository
}

// NewAccountService cria um novo serviço de contas
func NewAccountService(repository ports.AccountRepository) *AccountService {
	return &AccountService{repository: repository}
}

// CreateAccount cria uma nova conta e valida duplicidade de API Key
// Retorna ErrDuplicatedAPIKey se a chave já existir
func (s *AccountService) CreateAccount(newAccount *domain.Account) (*domain.Account, error) {
	// Verifica duplicidade de API Key antes da criação
	existingAccount, err := s.repository.FindByAPIKey(newAccount.ApiKey)

	if err != nil && err != domain.ErrAccountNotFound {
		return nil, err
	}

	if existingAccount != nil {
		return nil, domain.ErrDuplicatedAPIKey
	}

	err = s.repository.Save(newAccount)
	if err != nil {
		return nil, err
	}

	return newAccount, nil
}

// UpdateBalance atualiza o saldo de uma conta de forma thread-safe
// O amount pode ser positivo (crédito)
func (s *AccountService) UpdateBalance(apiKey string, amount float64) (*domain.Account, error) {
	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	account.AddBalance(amount)
	err = s.repository.UpdateBalance(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// FindByAPIKey busca uma conta pelo API Key
func (s *AccountService) FindByAPIKey(apiKey string) (*domain.Account, error) {
	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// FindByID busca uma conta pelo ID
func (s *AccountService) FindByID(id string) (*domain.Account, error) {
	account, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

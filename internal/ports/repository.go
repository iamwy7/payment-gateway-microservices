package ports

import "github.com/iamwy7/payment-gateway-management/internal/domain"

type AccountRepository interface {
	Save(account *domain.Account) error
	FindByAPIKey(apiKey string) (*domain.Account, error)
	FindByID(id string) (*domain.Account, error)
	UpdateBalance(account *domain.Account) error
}

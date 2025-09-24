package middleware

import (
	"net/http"

	"github.com/iamwy7/payment-gateway-management/internal/application/service"
	"github.com/iamwy7/payment-gateway-management/internal/domain"
)

type AuthMiddleware struct {
	accountService *service.AccountService
}

func NewAuthMiddleware(accountservice *service.AccountService) *AuthMiddleware {
	return &AuthMiddleware{
		accountService: accountservice,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			http.Error(w, "X-API-KEY is required", http.StatusUnauthorized)
			return
		}

		_, err := m.accountService.FindByAPIKey(apiKey)
		if err != nil {
			if err == domain.ErrAccountNotFound {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

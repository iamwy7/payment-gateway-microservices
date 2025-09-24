package domain

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

type PaymentType string

const (
	Card PaymentType = "card"
)

type InvoiceRequest struct {
	ApiKey         string
	Amount         float64
	Description    string
	PaymentType    PaymentType
	CardNumber     string
	CVV            string
	ExpiryMonth    int
	ExpiryYear     int
	CardholderName string
}

type Invoice struct {
	ID             string // TODO: Try UUID v7 to prevent performance issues with SQL databases
	AccountID      string
	Amount         float64
	Status         Status
	Description    string
	PaymentType    PaymentType
	CardLastDigits string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreditCardVO struct {
	Number         string
	CVV            string
	ExpiryMonth    int
	ExpiryYear     int
	CardholderName string
}

func NewInvoice(accountID string, amount float64, description string, paymentType PaymentType, card CreditCardVO) (*Invoice, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	lastDigits := card.Number[len(card.Number)-4:]

	return &Invoice{
		ID:             uuid.New().String(),
		AccountID:      accountID,
		Amount:         amount,
		Status:         StatusPending,
		Description:    description,
		PaymentType:    paymentType,
		CardLastDigits: lastDigits,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (i *Invoice) Process() error {
	if i.Amount >= 10000 {
		return nil
	}
	// That simulates an percentage of approval from BACEN
	randomSource := rand.New(rand.NewSource(time.Now().Unix()))
	var newStatus Status

	if randomSource.Float64() <= 0.7 {
		newStatus = StatusApproved
	} else {
		newStatus = StatusRejected
	}
	i.Status = newStatus
	return nil
}

func (i *Invoice) UpdateStatus(newStatus Status) error {
	if i.Status != StatusPending {
		return ErrInvalidStatus
	}
	i.Status = newStatus
	i.UpdatedAt = time.Now()
	return nil
}

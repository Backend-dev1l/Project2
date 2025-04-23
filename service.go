package tsService

import (
	"errors"
	"fmt"
)

type AccountingService interface {
	// Deposit - начисление средств на баланс пользователя
	Deposit(userID int64, amount float64) error

	// Reserve - резервирование средств с основного баланса
	Reserve(userID, orderID int64, amount float64) error

	// ConfirmRevenue - подтверждение выручки (списание из резерва)
	ConfirmRevenue(userID, orderID int64, amount float64) error

	// GetUserBalance - получение баланса пользователя
	GetBalance(userID int64) (float64, error)
}

type accountingService struct {
	repo AccountingRepository
}

func NewAccountingService(r AccountingRepository) AccountingService {
	return &accountingService{repo: r}
}

func (s *accountingService) ConfirmRevenue(userID int64, orderID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	// Подтверждаем резерв
	if err := s.repo.ConfirmRevenue(userID, orderID, amount); err != nil {
		return fmt.Errorf("failed to confirm revenue: %w", err)
	}
	return nil
}

func (s *accountingService) Deposit(userID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return s.repo.Deposit(userID, amount)

}

func (s *accountingService) GetBalance(userID int64) (float64, error) {
	balance, err := s.repo.GetBalance(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user balance: %w", err)
	}
	return balance, nil
}

func (s *accountingService) Reserve(userID int64, orderID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Получаем текущий баланс (основной + резерв)
	currentBalance, err := s.repo.GetBalance(userID)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	// Проверяем достаточно ли средств для резервирования
	if currentBalance < amount {
		return fmt.Errorf("insufficient funds: available %.2f, required %.2f", currentBalance, amount)
	}

	// Создаем резерв
	if err := s.repo.Reserve(userID, orderID, amount); err != nil {
		return fmt.Errorf("failed to create reserve: %w", err)
	}

	return nil
}

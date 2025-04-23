package tsService

import (
	"fmt"

	"gorm.io/gorm"
)

type AccountingRepository interface {
	// Метод начисления средств на баланс
	Deposit(userID int64, amount float64) error // post

	// Метод резервирования средств с основного баланса
	Reserve(userID int64, orderID int64, amount float64) error //post

	// Метод подтверждения выручки (списание из резерва)
	ConfirmRevenue(userID int64, orderID int64, amount float64) error //post

	// Метод получения баланса пользователя (основной баланс + резерв)
	GetBalance(userID int64) (float64, error) // get
}

type accountingRepository struct {
	db *gorm.DB
}

func NewAccountingRepository(db *gorm.DB) accountingRepository {
	return accountingRepository{db: db}
}

// Deposit начисляет средства на баланс
func (r *accountingRepository) Deposit(userID int64, amount float64) error {
	result := r.db.Exec("UPDATE accounts SET balance = balance + ? WHERE user_id = ?", amount, userID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *accountingRepository) Reserve(userID int64, orderID int64, amount float64) error {
	// Проверяем баланс пользователя
	var account struct{ Balance float64 }
	if err := r.db.Table("accounts").
		Select("balance").
		Where("user_id = ?", userID).
		First(&account).Error; err != nil {
		return fmt.Errorf("failed to get account balance: %w", err)
	}

	// Проверяем достаточно ли средств
	if account.Balance < amount {
		return fmt.Errorf("insufficient funds: available %.2f, requested %.2f",
			account.Balance, amount)
	}

	// Списываем средства с основного баланса
	result := r.db.Exec("UPDATE accounts SET balance = balance - ? WHERE user_id = ? AND balance >= ?", amount, userID, amount)
	if result.Error != nil {
		return fmt.Errorf("failed to deduct funds: %w", result.Error)
	}

	// Создаем запись о резерве
	reserve := struct {
		UserID  int64   `gorm:"column:user_id"`
		OrderID int64   `gorm:"column:order_id"`
		Amount  float64 `gorm:"column:amount"`
	}{
		UserID:  userID,
		OrderID: orderID,
		Amount:  amount,
	}

	if err := r.db.Table("reserves").Create(&reserve).Error; err != nil {
		return fmt.Errorf("failed to create reserve record: %w", err)
	}

	return nil
}

// ConfirmRevenue подтверждает выручку
func (r *accountingRepository) ConfirmRevenue(userID int64, orderID int64, amount float64) error {
	// Проверяем резерв
	var reserve struct{ Amount float64 }
	if err := r.db.Table("reserves").
		Select("amount").
		Where("user_id = ? AND order_id = ?", userID, orderID).First(&reserve).Error; err != nil {
		return fmt.Errorf("failed to check reserve: %w", err)
	}

	//  Проверяем достаточно ли зарезервированных средств
	if reserve.Amount < amount {
		return fmt.Errorf("insufficient reserved amount: available %.2f, requested %.2f",
			reserve.Amount, amount)
	}

	// уменьшаем резерв
	result := r.db.Exec(
		"UPDATE reserves SET amount = amount - ? WHERE user_id = ? AND order_id = ? AND amount >= ?",
		amount, userID, orderID, amount)

	if result.Error != nil {
		return fmt.Errorf("failed to update reserve: %w", result.Error)
	}

	// 4. Записываем в историю
	history := struct {
		UserID  int64   `gorm:"column:user_id"`
		OrderID int64   `gorm:"column:order_id"`
		Amount  float64 `gorm:"column:amount"`
	}{
		UserID:  userID,
		OrderID: orderID,
		Amount:  amount,
	}

	if err := r.db.Table("revenue_history").Create(&history).Error; err != nil {
		return fmt.Errorf("failed to record revenue history: %w", err)
	}

	return nil
}

func (r *accountingRepository) GetBalance(userID int64) (float64, error) {
	// Получаем основной баланс
	var account struct {
		Balance float64
	}

	if err := r.db.Table("accounts").
		Select("balance").
		Where("user_id = ?", userID).
		First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("account not found")
		}
		return 0, fmt.Errorf("failed to get account balance: %w", err)
	}

	// Получаем сумму зарезервированных средств
	var reservedAmount float64
	if err := r.db.Table("reserves").
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ?", userID).
		Scan(&reservedAmount).Error; err != nil {
		return 0, fmt.Errorf("failed to get reserved amount: %w", err)
	}

	// Возвращаем общий баланс
	return account.Balance + reservedAmount, nil
}

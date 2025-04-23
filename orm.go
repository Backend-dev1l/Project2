package tsService

type Account struct {
	UserID  int64 `gorm:"primaryKey"`
	Balance float64
}

type Reserve struct {
	ID        uint `gorm:"primaryKey"`
	UserID    int64
	ServiceID int64
	OrderID   int64
	Amount    float64
}

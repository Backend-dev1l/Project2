package handlerstructs

type DepositRequest struct {
	// Пополнение (Deposit)
	user_id int64 `json:"user_id"`
	amount  int64 `json:"amount"` // в копейках
}

type DepositResponse struct {
	message string `json:"message"`
}

type ReserveRequest struct {
	// Резервирование средств (Reserve)
	user_id  int64 `json:"user_id"`
	order_id int64 `json:"order_id"`
	amount   int64 `json:"amount"` // в копейках
}

type ReserveResponse struct {
	status string `json:"status"`
}

type ConfirmRevenueRequest struct {
	// Подтверждение выручки (ConfirmRevenue)
	user_id  int64 `json:"user_id"`
	order_id int64 `json:"order_id"`
	amount   int64 `json:"amount"` // в копейках
}

type ConfirmRevenueResponse struct {
	status string `json:"status"`
}

type GetBalanceRequest struct {
	// Получение текущего баланса
	user_id int64 `json:"user_id"`
}

type GetBalanceResponse struct {
	user_id int64 `json:"user_id"`
	amount  int64 `json:"amount"`
}

type ErrorResponse struct {
	Error string `json:"error`
}

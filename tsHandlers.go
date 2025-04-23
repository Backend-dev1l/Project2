package handlers

import (
	tsService "RestApi/internal/services/ts"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AccountingHandler struct {
	service tsService.AccountingService
}

func NewAccountingHandler(s tsService.AccountingService) *AccountingHandler {
	return &AccountingHandler{service: s}
}

func (h *AccountingHandler) GetBalanceHandler(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if userIDStr == "" || err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid ID",
		})
	}

	balance, err := h.service.GetBalance(int64(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Can't get balance",
		})
	}

	return c.JSON(http.StatusOK, GetBalanceResponse{
		user_id: int64(userID),
		amount:  int64(balance),
	})
}

func (h *AccountingHandler) ConfirmRevenueHandler(c echo.Context) error {
	var revenue ConfirmRevenueRequest
	if err := c.Bind(&revenue); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
	}

	if revenue.user_id == 0 || revenue.order_id == 0 || revenue.amount == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Missing required parameters"})
	}

	err := h.service.ConfirmRevenue(revenue.user_id, revenue.order_id, float64(revenue.amount))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "revenue confirmation failed",
		})
	}
	return c.JSON(http.StatusOK, ConfirmRevenueResponse{
		status: "Success",
	})
}

func (h *AccountingHandler) ReserveHandler(c echo.Context) error {
	var reserve ReserveRequest
	if err := c.Bind(&reserve); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request format",
		})
	}

	if reserve.user_id == 0 || reserve.order_id == 0 || reserve.amount == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Missing required parameters"})
	}

	err := h.service.Reserve(reserve.user_id, reserve.order_id, float64(reserve.amount))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "reservation failed",
		})
	}
	return c.JSON(http.StatusOK, ReserveResponse{
		status: "Success",
	})
}

func (h *AccountingHandler) DepositHandler(c echo.Context) error {
	var deposit DepositRequest
	if err := c.Bind(&deposit); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request format",
		})
	}

	if deposit.user_id == 0 || deposit.amount == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Missing required parameters"})
	}

	err := h.service.Deposit(deposit.user_id, float64(deposit.amount))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "deposit failed",
		})
	}
	return c.JSON(http.StatusOK, DepositResponse{
		message: "balance successfully updated",
	})
}

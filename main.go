package main

import (
	"RestApi/internal/db"
	tsService "RestApi/internal/services/ts"

	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	e := echo.New()

	accountRepo := tsService.NewAccountingRepository(database)
	accountService := tsService.NewAccountingService(&accountRepo)
	accountHandlers := handlers.NewAccountingHandler(accountService)

	e.POST("/Reserve", accountHandlers.ReserveHandler)
	e.POST("/Deposit", accountHandlers.DepositHandler)
	e.POST("/ConfirmRevenue", accountHandlers.ConfirmRevenueHandler)
	e.GET("/GetBalance", accountHandlers.GetBalanceHandler)
	if err := e.Start(":8083"); err != nil {
		log.Fatal(err)
	}
}

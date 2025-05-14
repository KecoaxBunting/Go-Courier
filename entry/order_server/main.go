package main

import (
	"go-courier/order"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env")
	db := order.InitDatabase()
	service := order.NewOrderService(db)
	order.StartGRPCServer(8090, service)
}

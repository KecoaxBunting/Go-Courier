package main

import (
	"go-courier/order"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env")
	order.InitDatabase()
	service := order.NewOrderService()
	order.StartGRPCServer(8090, service)
}

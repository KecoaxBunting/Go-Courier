package main

import (
	"go-courier/courier"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env")
	db := courier.InitDatabase()
	service := courier.NewCourierService(db)
	courier.StartGRPCServer(8100, service)
}

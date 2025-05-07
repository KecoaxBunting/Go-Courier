package main

import (
	"go-courier/courier"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env")
	courier.InitDatabase()
	service := courier.NewCourierService()
	courier.StartGRPCServer(8100, service)
}

package main

import (
	"go-courier/auth"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env")
	auth.InitDatabase()
	service := auth.NewAuthService()
	auth.StartGRPCServer(8080, service)
}

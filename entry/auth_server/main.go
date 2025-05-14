package main

import (
	"go-courier/auth"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env")
	db := auth.InitDatabase()
	service := auth.NewAuthService(db)
	auth.StartGRPCServer(8080, service)
}

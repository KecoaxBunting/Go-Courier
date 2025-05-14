package main

import (
	"fmt"
	"go-courier/delivery"
	"go-courier/proto/courier"
	"go-courier/proto/order"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	_ = godotenv.Load("../../.env")
	db := delivery.InitDatabase()

	courierHost := fmt.Sprintf("%s:8100", os.Getenv("COURIER_SERVICE_HOST"))
	courierConn, err := grpc.NewClient(courierHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to courier service: %v", err)
	}
	defer courierConn.Close()

	courierClient := courier.NewCourierServiceClient(courierConn)

	orderHost := fmt.Sprintf("%s:8090", os.Getenv("ORDER_SERVICE_HOST"))
	orderConn, err := grpc.NewClient(orderHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer orderConn.Close()

	orderClient := order.NewOrderServiceClient(orderConn)

	service := delivery.NewDeliveryService(db, courierClient, orderClient)
	delivery.StartGRPCServer(8110, service)
}

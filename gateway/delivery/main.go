package main

import (
	"go-courier/controller"
	"go-courier/middleware"
	deliverypb "go-courier/proto/delivery"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	_ = godotenv.Load("../../.env")
	conn, err := grpc.NewClient(os.Getenv("DELIVERY_SERVICE_HOST")+":8110", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to gRPC server :%v", err)
	}
	defer conn.Close()

	client := deliverypb.NewDeliveryServiceClient(conn)

	router := gin.Default()

	router.Use(middleware.JWTAuth())
	{
		router.POST("/assignCourier", controller.AssignCourier(client))

		router.PUT("/completeOrder/:deliveryId", controller.CompleteOrder(client))

		router.GET("/deliveries/:deliveryId", controller.GetDelivery(client))

		router.GET("/deliveries/", controller.ListDelivery(client))

		router.PUT("/deliveries/:deliveryId", controller.UpdateDelivery(client))

		router.DELETE("/deliveries/:deliveryId", controller.DeleteDelivery(client))

	}
	router.Run(":8111")
}

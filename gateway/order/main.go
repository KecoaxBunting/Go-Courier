package main

import (
	"go-courier/controller"
	"go-courier/middleware"
	orderpb "go-courier/proto/order"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	_ = godotenv.Load("../../.env")
	conn, err := grpc.NewClient(os.Getenv("ORDER_SERVICE_HOST")+":8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to gRPC server :%v", err)
	}
	defer conn.Close()

	client := orderpb.NewOrderServiceClient(conn)

	router := gin.Default()

	protected := router.Group("/orders")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/", controller.CreateOrder(client))

		protected.GET("/:orderId", controller.GetOrder(client))

		protected.GET("/", controller.ListOrder(client))

	}
	router.Run(":8091")
}

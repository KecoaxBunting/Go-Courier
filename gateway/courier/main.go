package main

import (
	"go-courier/controller"
	"go-courier/middleware"
	courierpb "go-courier/proto/courier"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	_ = godotenv.Load("../../.env")
	conn, err := grpc.NewClient(os.Getenv("COURIER_SERVICE_HOST")+":8100", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to gRPC server :%v", err)
	}
	defer conn.Close()

	client := courierpb.NewCourierServiceClient(conn)

	router := gin.Default()

	protected := router.Group("/couriers")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/", controller.RegisterCourier(client))

		protected.GET("/:courierId", controller.GetCourier(client))

		protected.GET("/", controller.ListCouriers(client))

	}
	router.Run(":8101")
}

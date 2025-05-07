package main

import (
	"go-courier/controller"
	"go-courier/middleware"
	authpb "go-courier/proto/auth"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	_ = godotenv.Load("../../.env")
	conn, err := grpc.NewClient(os.Getenv("AUTH_SERVICE_HOST")+":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to gRPC server :%v", err)
	}
	defer conn.Close()

	client := authpb.NewAuthServiceClient(conn)

	router := gin.Default()

	router.POST("/register", controller.Register(client))

	router.POST("/login", controller.Login(client))

	protected := router.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/me", func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(200, gin.H{"user": user})
		})
	}
	router.Run(":8081")
}

package controller

import (
	"context"
	authpb "go-courier/proto/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(client authpb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req = authpb.RegisterRequest{}

		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		res, err := client.Register(context.Background(), &authpb.RegisterRequest{
			Username: req.Username,
			Password: req.Password,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func Login(client authpb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req authpb.LoginRequest

		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		res, err := client.Login(context.Background(), &authpb.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

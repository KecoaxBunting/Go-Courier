package controller

import (
	"context"
	"errors"
	helper "go-courier/helper"
	interceptor "go-courier/interceptor"
	orderpb "go-courier/proto/order"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrder(client orderpb.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req = orderpb.OrderRequest{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		userIdRaw, ok := c.Get(string(interceptor.UserIdKey))
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
			return
		}

		userId, ok := userIdRaw.(int64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid user Id"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		ctx := helper.GRPCWithAuth(context.Background(), strings.TrimPrefix(authHeader, "Bearer "))

		res, err := client.CreateOrder(ctx, &orderpb.OrderRequest{
			SenderId: userId,
			Items:    req.Items,
			Address:  req.Address,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func GetOrder(client orderpb.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("orderId")
		id, err := strconv.ParseInt(param, 0, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		ctx := helper.GRPCWithAuth(context.Background(), strings.TrimPrefix(authHeader, "Bearer "))

		res, err := client.GetOrder(ctx, &orderpb.GetOrderRequest{Id: id})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func ListOrder(client orderpb.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		ctx := helper.GRPCWithAuth(context.Background(), strings.TrimPrefix(authHeader, "Bearer "))

		res, err := client.ListOrder(ctx, &orderpb.Empty{})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

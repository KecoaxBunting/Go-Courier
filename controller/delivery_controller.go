package controller

import (
	"context"
	"errors"
	helper "go-courier/helper"
	deliverypb "go-courier/proto/delivery"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AssignCourier(client deliverypb.DeliveryServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req = deliverypb.DeliveryRequest{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		} else if req.CourierId == nil || req.OrderId == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Input is required"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		ctx := helper.GRPCWithAuth(context.Background(), strings.TrimPrefix(authHeader, "Bearer "))

		res, err := client.AssignCourier(ctx, &deliverypb.DeliveryRequest{
			CourierId: req.CourierId,
			OrderId:   req.OrderId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func CompleteOrder(client deliverypb.DeliveryServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("deliveryId")
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

		res, err := client.CompleteOrder(ctx, &deliverypb.CompleteRequest{Id: id})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Delivery data not found"})
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func GetDelivery(client deliverypb.DeliveryServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("deliveryId")
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

		res, err := client.GetDelivery(ctx, &deliverypb.GetDeliveryRequest{Id: id})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func ListDelivery(client deliverypb.DeliveryServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		ctx := helper.GRPCWithAuth(context.Background(), strings.TrimPrefix(authHeader, "Bearer "))

		res, err := client.ListDelivery(ctx, &deliverypb.EmptyDelivery{})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func DeleteDelivery(client deliverypb.DeliveryServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("deliveryId")
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
		res, err := client.DeleteDelivery(ctx, &deliverypb.DeleteDeliveryRequest{Id: id})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

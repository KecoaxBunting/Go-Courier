package controller

import (
	"context"
	"errors"
	helper "go-courier/helper"
	courierpb "go-courier/proto/courier"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCourier(client courierpb.CourierServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req = courierpb.CourierRequest{}
		phoneRegex := regexp.MustCompile(`^(?:\+62|62|0)8[1-9][0-9]{6,9}$`)
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		} else if req.Name == nil || req.PhoneNumber == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Input is required"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		ctx := helper.GRPCWithAuth(context.Background(), strings.TrimPrefix(authHeader, "Bearer "))

		match := phoneRegex.MatchString(*req.PhoneNumber)

		if !match {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone number format"})
			return
		}

		res, err := client.RegisterCourier(ctx, &courierpb.CourierRequest{
			Name:        req.Name,
			PhoneNumber: req.PhoneNumber,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func GetCourier(client courierpb.CourierServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("courierId")
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

		res, err := client.GetCourier(ctx, &courierpb.GetCourierRequest{Id: id})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func ListCouriers(client courierpb.CourierServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		ctx := helper.GRPCWithAuth(context.Background(), strings.TrimPrefix(authHeader, "Bearer "))

		res, err := client.ListCourier(ctx, &courierpb.EmptyCourier{})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func UpdateCourier(client courierpb.CourierServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("courierId")
		id, err := strconv.ParseInt(param, 0, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
			return
		}

		var req = courierpb.UpdateCourierRequest{}
		err = c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		} else if req.Name == nil || req.PhoneNumber == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Input is required"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		ctx := helper.GRPCWithAuth(context.Background(), strings.TrimPrefix(authHeader, "Bearer "))
		res, err := client.UpdateCourier(ctx, &courierpb.UpdateCourierRequest{
			Id:          id,
			Name:        req.Name,
			PhoneNumber: req.PhoneNumber,
		})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func DeleteCourier(client courierpb.CourierServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("courierId")
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
		res, err := client.DeleteCourier(ctx, &courierpb.DeleteCourierRequest{Id: id})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func ChangeAvailability(client courierpb.CourierServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("courierId")
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
		res, err := client.ChangeAvailability(ctx, &courierpb.ChangeAvailabilityCourierRequest{Id: id})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

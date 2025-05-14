package delivery

import (
	"context"
	"errors"
	"fmt"
	"go-courier/helper"
	interceptor "go-courier/interceptor"
	courierpb "go-courier/proto/courier"
	deliverypb "go-courier/proto/delivery"
	orderpb "go-courier/proto/order"

	"gorm.io/gorm"
)

type DeliveryService struct {
	deliverypb.UnimplementedDeliveryServiceServer
	courierpb.CourierServiceClient
	orderpb.OrderServiceClient
	*gorm.DB
}

func NewDeliveryService(db *gorm.DB, courierClient courierpb.CourierServiceClient, orderClient orderpb.OrderServiceClient) *DeliveryService {
	return &DeliveryService{
		DB:                   db,
		CourierServiceClient: courierClient,
		OrderServiceClient:   orderClient,
	}
}

func (d *DeliveryService) AssignCourier(ctx context.Context, req *deliverypb.DeliveryRequest) (*deliverypb.DeliveryResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}

	newCtx := helper.ForwardMetadata(ctx)

	courier, err := d.CourierServiceClient.GetCourier(newCtx, &courierpb.GetCourierRequest{Id: *req.CourierId})
	if err != nil {
		return nil, err
	} else if !courier.Available {
		return nil, fmt.Errorf("courier is not available on this time")
	}

	order, err := d.OrderServiceClient.GetOrder(newCtx, &orderpb.GetOrderRequest{Id: *req.OrderId})
	if err != nil {
		return nil, err
	} else if order.Status == "On Delivering" {
		return nil, fmt.Errorf("can no make new asssignment, order is already on delivering")
	} else if order.Status == "Complete" {
		return nil, fmt.Errorf("order is completed")
	}

	var existDelivery = &Delivery{}
	err = d.DB.Where("courier_id = ? AND order_id = ?", req.CourierId, req.OrderId).First(&existDelivery).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		var delivery = &Delivery{
			CourierId: *req.CourierId,
			OrderId:   *req.OrderId,
			Status:    "On Progress",
			AddedBy:   userId,
		}

		err := d.DB.Create(&delivery).Error
		if err != nil {
			return nil, fmt.Errorf("failed to assign courier")
		}

		_, err = d.OrderServiceClient.SetOrderToDelivering(newCtx, &orderpb.SetOrderRequest{Id: delivery.OrderId})
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check exist data")
	} else {
		if existDelivery.Status == "Complete" {
			return nil, fmt.Errorf("data exist and has been completed")
		} else if existDelivery.Status == "On Progress" {
			return nil, fmt.Errorf("data exist")
		}
	}

	return &deliverypb.DeliveryResponse{
		Message: "Successfully assign driver to order",
		Status:  "On Progress",
	}, nil
}

func (d *DeliveryService) CompleteOrder(ctx context.Context, req *deliverypb.CompleteRequest) (*deliverypb.DeliveryResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}

	var delivery = Delivery{}
	err := d.DB.Where("id = ? AND added_by = ?", req.Id, userId).First(&delivery).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("no delivery data found")
	} else if err != nil {
		return nil, fmt.Errorf("can not get delivery data")
	} else if delivery.Status == "Complete" {
		return nil, fmt.Errorf("delivery data has been completed")
	}

	delivery.Status = "Complete"
	err = d.DB.Model(&Delivery{}).Where("id = ?", delivery.Id).Update("status", delivery.Status).Error
	if err != nil {
		return nil, fmt.Errorf("failed to complete delivery")
	}

	newCtx := helper.ForwardMetadata(ctx)
	_, err = d.OrderServiceClient.SetOrderToComplete(newCtx, &orderpb.SetOrderRequest{Id: delivery.OrderId})
	if err != nil {
		return nil, err
	}

	return &deliverypb.DeliveryResponse{
		Message: "Successfully complete delivery and order",
		Status:  delivery.Status,
	}, nil
}

func (d *DeliveryService) GetDelivery(ctx context.Context, req *deliverypb.GetDeliveryRequest) (*deliverypb.Delivery, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}

	var delivery = Delivery{}
	err := d.DB.Where("id = ? AND added_by = ?", req.Id, userId).First(&delivery).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("no delivery data found")
	} else if err != nil {
		return nil, fmt.Errorf("can not get delivery data")
	}

	newCtx := helper.ForwardMetadata(ctx)

	courier, err := d.CourierServiceClient.GetCourier(newCtx, &courierpb.GetCourierRequest{Id: delivery.CourierId})
	if err != nil {
		return nil, err
	}

	order, err := d.OrderServiceClient.GetOrder(newCtx, &orderpb.GetOrderRequest{Id: delivery.OrderId})
	if err != nil {
		return nil, err
	}

	return &deliverypb.Delivery{
		Id:          delivery.Id,
		CourierData: courier,
		OrderData:   order,
		Status:      delivery.Status,
		AddedBy:     courier.AddedBy,
	}, nil
}

func (d *DeliveryService) ListDelivery(ctx context.Context, req *deliverypb.EmptyDelivery) (*deliverypb.DeliveryList, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}

	var deliveries = []Delivery{}
	err := d.DB.Where("added_by = ?", userId).Find(&deliveries).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("no deliveries data found")
	} else if err != nil {
		return nil, fmt.Errorf("can not get deliveries data")
	}

	var deliveryResponse = []*deliverypb.Delivery{}
	newCtx := helper.ForwardMetadata(ctx)

	for _, value := range deliveries {
		courier, err := d.CourierServiceClient.GetCourier(newCtx, &courierpb.GetCourierRequest{Id: value.CourierId})
		if err != nil {
			return nil, err
		}

		order, err := d.OrderServiceClient.GetOrder(newCtx, &orderpb.GetOrderRequest{Id: value.OrderId})
		if err != nil {
			return nil, err
		}

		deliveryResponse = append(deliveryResponse, &deliverypb.Delivery{
			Id:          value.Id,
			CourierData: courier,
			OrderData:   order,
			Status:      value.Status,
			AddedBy:     value.AddedBy,
		})
	}

	return &deliverypb.DeliveryList{
		Deliveries: deliveryResponse,
	}, nil
}

func (d *DeliveryService) UpdateDelivery(ctx context.Context, req *deliverypb.UpdateDeliveryRequest) (*deliverypb.DeliveryResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}

	var delivery = Delivery{}
	err := d.DB.Where("id = ? AND added_by = ?", req.Id, userId).First(&delivery).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("delivery data not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get delivery data")
	} else if delivery.Status == "Complete" {
		return nil, fmt.Errorf("can not update completed delivery")
	}

	newCtx := helper.ForwardMetadata(ctx)

	_, err = d.CourierServiceClient.GetCourier(newCtx, &courierpb.GetCourierRequest{Id: delivery.CourierId})
	if err != nil {
		return nil, err
	}

	_, err = d.OrderServiceClient.GetOrder(newCtx, &orderpb.GetOrderRequest{Id: delivery.OrderId})
	if err != nil {
		return nil, err
	}

	delivery.CourierId = *req.CourierId
	delivery.OrderId = *req.OrderId

	err = d.DB.Updates(&delivery).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update delivery data")
	}

	return &deliverypb.DeliveryResponse{
		Message: "Successfully update delivery data",
	}, nil
}

func (d *DeliveryService) DeleteDelivery(ctx context.Context, req *deliverypb.DeleteDeliveryRequest) (*deliverypb.DeliveryResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}

	err := d.DB.Where("id = ? AND added_by = ?", req.Id, userId).Delete(&Delivery{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("order not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to delete order")
	}

	return &deliverypb.DeliveryResponse{
		Message: "Successfully delete order",
	}, nil
}

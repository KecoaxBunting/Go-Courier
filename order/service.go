package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	interceptor "go-courier/interceptor"
	orderpb "go-courier/proto/order"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	orderpb.UnimplementedOrderServiceServer
	*gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		DB: db,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, req *orderpb.OrderRequest) (*orderpb.OrderResponse, error) {
	itemsJSON, err := json.Marshal(req.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to encode JSON")
	}

	var order = Order{
		SenderId: req.SenderId,
		Items:    itemsJSON,
		Address:  *req.Address,
		Status:   "Pending",
	}

	err = o.DB.Create(&order).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create new order")
	}

	return &orderpb.OrderResponse{
		Id:     order.Id,
		Status: order.Status,
	}, nil
}

func (o *OrderService) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}
	var order = Order{}
	err := o.DB.Where("id = ? AND sender_id = ?", req.Id, userId).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("order not found")
	} else if err != nil {
		return nil, fmt.Errorf("can not get order")
	}

	var items = []string{}
	err = json.Unmarshal(order.Items, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON")
	}

	return &orderpb.Order{
		Id:        order.Id,
		SenderId:  order.SenderId,
		Items:     items,
		Address:   order.Address,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (o *OrderService) ListOrder(ctx context.Context, _ *orderpb.EmptyOrder) (*orderpb.OrderList, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}
	var orders = []Order{}
	err := o.DB.Where("sender_id = ?", userId).Find(&orders).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("no orders found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get orders")
	}

	var listOrder = []*orderpb.Order{}
	for _, value := range orders {
		var items = []string{}
		err := json.Unmarshal(value.Items, &items)
		if err != nil {
			return nil, fmt.Errorf("failed to decode JSON")
		}

		listOrder = append(listOrder, &orderpb.Order{
			Id:        value.Id,
			SenderId:  value.SenderId,
			Items:     items,
			Address:   value.Address,
			Status:    value.Status,
			CreatedAt: value.CreatedAt.Format(time.RFC3339),
		})
	}
	return &orderpb.OrderList{Orders: listOrder}, nil
}

func (o *OrderService) UpdateOrder(ctx context.Context, req *orderpb.UpdateOrderRequest) (*orderpb.UpdateOrderResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}
	var order = Order{}
	err := o.DB.Where("id = ? AND sender_id = ?", req.Id, userId).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("order not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get order")
	} else if order.Status == "Complete" {
		return nil, fmt.Errorf("can not update completed order")
	} else if order.Status == "On Delivering" {
		return nil, fmt.Errorf("can not update on delivering order")
	}

	itemsJSON, err := json.Marshal(req.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to encode JSON")
	}

	order.Items = itemsJSON
	order.Address = *req.Address

	err = o.DB.Updates(&order).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update order")
	}

	return &orderpb.UpdateOrderResponse{
		Message: "Successfully update order",
	}, nil
}

func (o *OrderService) DeleteOrder(ctx context.Context, req *orderpb.DeleteOrderRequest) (*orderpb.DeleteOrderResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}

	var order = Order{}
	err := o.DB.Where("id = ? AND sender_id = ?", req.Id, userId).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("order not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get order")
	} else if order.Status == "On Delivering" {
		return nil, fmt.Errorf("can not delete on delivering order")
	}

	err = o.DB.Delete(&order).Error
	if err != nil {
		return nil, fmt.Errorf("failed to delete order")
	}

	return &orderpb.DeleteOrderResponse{
		Message: "Successfully delete order",
	}, nil
}

func (o *OrderService) SetOrderToComplete(ctx context.Context, req *orderpb.SetOrderRequest) (*orderpb.SetOrderResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}

	var orderData = Order{}
	err := o.DB.Where("id = ? AND sender_id = ?", req.Id, userId).First(&orderData).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("order not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get order")
	}

	orderData.Status = "Complete"
	err = o.DB.Model(&Order{}).Where("id = ?", orderData.Id).Update("status", orderData.Status).Error
	if err != nil {
		return nil, fmt.Errorf("failed to set order to complete")
	}

	return &orderpb.SetOrderResponse{
		Message: "Successfully set order to complete",
	}, nil
}

func (o *OrderService) SetOrderToDelivering(ctx context.Context, req *orderpb.SetOrderRequest) (*orderpb.SetOrderResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}

	var orderData = Order{}
	err := o.DB.Where("id = ? AND sender_id = ?", req.Id, userId).First(&orderData).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("order not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get order")
	}

	orderData.Status = "On Delivering"
	err = o.DB.Model(&Order{}).Where("id = ?", orderData.Id).Update("status", orderData.Status).Error
	if err != nil {
		return nil, fmt.Errorf("failed to set order to delivering")
	}

	return &orderpb.SetOrderResponse{
		Message: "Successfully set order to delivering",
	}, nil
}

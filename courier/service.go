package courier

import (
	"context"
	"errors"
	"fmt"
	interceptor "go-courier/interceptor"
	courierpb "go-courier/proto/courier"

	"gorm.io/gorm"
)

type CourierService struct {
	courierpb.UnimplementedCourierServiceServer
}

func NewCourierService() *CourierService {
	return &CourierService{}
}

func (c *CourierService) RegisterCourier(ctx context.Context, req *courierpb.CourierRequest) (*courierpb.CourierResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}
	courier := &Courier{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Available:   true,
		AddedBy:     userId,
	}
	err := DB.Create(&courier).Error
	if err != nil {
		return nil, fmt.Errorf("failed to register new courier")
	}
	return &courierpb.CourierResponse{
		Message: "Successfully register new courier",
	}, nil
}

func (c *CourierService) GetCourier(ctx context.Context, req *courierpb.GetCourierRequest) (*courierpb.Courier, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}
	var courier = Courier{}
	err := DB.Where("id = ? AND added_by = ?", req.Id, userId).First(&courier).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("courier not found")
	} else if err != nil {
		return nil, fmt.Errorf("can not get courier")
	}

	return &courierpb.Courier{
		Id:          courier.Id,
		Name:        courier.Name,
		PhoneNumber: courier.PhoneNumber,
		Available:   courier.Available,
		AddedBy:     courier.AddedBy,
	}, nil
}

func (c *CourierService) ListCourier(ctx context.Context, _ *courierpb.Empty) (*courierpb.CourierList, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}

	var couriers = []Courier{}
	err := DB.Where("added_by = ?", userId).Find(&couriers).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("no couriers found")
	} else if err != nil {
		return nil, fmt.Errorf("can not get all couriers")
	}

	var CourierResponse = []*courierpb.Courier{}
	for _, value := range couriers {
		CourierResponse = append(CourierResponse, &courierpb.Courier{
			Id:          value.Id,
			Name:        value.Name,
			PhoneNumber: value.PhoneNumber,
			Available:   value.Available,
			AddedBy:     value.AddedBy,
		})
	}
	return &courierpb.CourierList{Couriers: CourierResponse}, nil
}

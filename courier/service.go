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
	*gorm.DB
}

func NewCourierService(db *gorm.DB) *CourierService {
	return &CourierService{
		DB: db,
	}
}

func (c *CourierService) RegisterCourier(ctx context.Context, req *courierpb.CourierRequest) (*courierpb.CourierResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}
	var courier = &Courier{
		Name:        *req.Name,
		PhoneNumber: *req.PhoneNumber,
		Available:   true,
		AddedBy:     userId,
	}
	err := c.DB.Create(&courier).Error
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
	err := c.DB.Where("id = ? AND added_by = ?", req.Id, userId).First(&courier).Error

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

func (c *CourierService) ListCourier(ctx context.Context, _ *courierpb.EmptyCourier) (*courierpb.CourierList, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}

	var couriers = []Courier{}
	err := c.DB.Where("added_by = ?", userId).Find(&couriers).Error
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

func (c *CourierService) UpdateCourier(ctx context.Context, req *courierpb.UpdateCourierRequest) (*courierpb.CourierResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}
	var courier = Courier{}
	err := c.DB.Where("id = ? AND added_by = ?", req.Id, userId).First(&courier).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("courier not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get courier")
	}

	courier.Name = *req.Name
	courier.PhoneNumber = *req.PhoneNumber

	err = c.DB.Updates(&courier).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update courier")
	}

	return &courierpb.CourierResponse{
		Message: "Successfully update courier",
	}, nil
}

func (c *CourierService) DeleteCourier(ctx context.Context, req *courierpb.DeleteCourierRequest) (*courierpb.CourierResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}
	err := c.DB.Where("id = ? AND added_by = ?", req.Id, userId).Delete(&Courier{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("courier not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to delete courier")
	}
	return &courierpb.CourierResponse{
		Message: "Successfully delete courier",
	}, nil
}

func (c *CourierService) ChangeAvailability(ctx context.Context, req *courierpb.ChangeAvailabilityCourierRequest) (*courierpb.CourierResponse, error) {
	userId, ok := ctx.Value(interceptor.UserIdKey).(int64)
	if !ok {
		return nil, fmt.Errorf("please login first")
	}
	var courier = Courier{}
	err := c.DB.Where("id = ? AND added_by = ?", req.Id, userId).First(&courier).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("courier not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get courier")
	}

	courier.Available = !courier.Available
	action := "disable"
	if courier.Available {
		action = "enable"
	}
	err = c.DB.Model(&Courier{}).Where("id = ?", courier.Id).Update("available", courier.Available).Error
	if err != nil {
		return nil, fmt.Errorf("failed to %s courier", action)
	}

	return &courierpb.CourierResponse{
		Message: fmt.Sprintf("Successfully %s courier", action),
	}, nil
}

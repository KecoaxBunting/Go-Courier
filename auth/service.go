package auth

import (
	"context"
	"fmt"
	helper "go-courier/helper"
	authpb "go-courier/proto/auth"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (a *AuthService) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	var user = User{}
	err := DB.Where("username = ?", req.Username).First(&user).Error
	if err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	var newUser = User{
		Username: req.Username,
		Password: string(hashed),
	}

	err = DB.Create(&newUser).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	return &authpb.AuthResponse{
		Token:   "",
		Message: "Register Success",
	}, nil
}

func (a *AuthService) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	var user = User{}
	err := DB.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := helper.GenerateToken(user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	return &authpb.AuthResponse{
		Token:   token,
		Message: "Login Success",
	}, nil
}

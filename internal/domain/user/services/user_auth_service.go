package services

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"user-service/internal/domain/user/dtos"
	"user-service/internal/domain/user/models"
	"user-service/pkg/auth"
	"user-service/pkg/exception"
	"user-service/pkg/hash"
	"user-service/pkg/logger"
)

func (service serviceUser) Register(ctx context.Context, input dtos.RegisterDto) (*dtos.CreateUserResponseDto, error) {
	emailExist, _ := service.Repository.FindByEmail(ctx, input.Email)
	fmt.Println("emaill: ", emailExist)
	if emailExist != nil {
		return nil, &exception.ErrWithCode{
			Code: http.StatusBadRequest,
			Err:  errors.New("email is already used"),
		}
	}

	if input.PhoneNumber != nil {
		phoneExist, _ := service.Repository.FindByPhone(ctx, *input.PhoneNumber)
		if phoneExist != nil {
			return nil, &exception.ErrWithCode{
				Code: http.StatusBadRequest,
				Err:  errors.New("phone is already used"),
			}
		}
	}

	password, err := hash.HashingPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user, err := service.Repository.Create(ctx, models.User{
		Password:    password,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
		IsActive:    false,
	})
	if err != nil {
		logger.Error("create user errors: ", zap.Error(err))
		return nil, err
	}

	return &dtos.CreateUserResponseDto{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}, nil
}

func (service serviceUser) Login(ctx context.Context, input dtos.LoginDto) (*dtos.TokenDto, error) {
	user, _ := service.Repository.FindByEmail(ctx, input.Email)
	if user == nil {
		return nil, &exception.ErrWithCode{
			Code: http.StatusNotFound,
			Err:  errors.New("user not found"),
		}
	}

	if !user.IsActive {
		return nil, &exception.ErrWithCode{
			Code: http.StatusBadRequest,
			Err:  errors.New("user inactive"),
		}
	}

	err := hash.ComparePassword(user.Password, input.Password)
	if err != nil {
		return nil, &exception.ErrWithCode{
			Code: http.StatusUnauthorized,
			Err:  errors.New("invalid password"),
		}
	}

	token, err := auth.CreateToken(auth.JWTPayload{
		UserID:    user.ID,
		FirstName: user.FirstName,
		Email:     user.Email,
		Role:      string(user.Role),
	})

	return &dtos.TokenDto{
		Token: token,
		Type:  "ACCESS_TOKEN",
	}, nil
}

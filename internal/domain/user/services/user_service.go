package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
	"user-service/internal/domain/user/dtos"
	"user-service/internal/domain/user/models"
	"user-service/internal/domain/user/repositories"
	"user-service/internal/shared"
	"user-service/pkg/exception"
	"user-service/pkg/hash"
	"user-service/pkg/helper"
)

func NewServiceUser(repository repositories.UserRepository, redis redis.Client) *serviceUser {
	return &serviceUser{
		Repository: repository,
		Redis:      redis,
	}
}

func (service serviceUser) Create(ctx context.Context, input dtos.CreateUserDto) (*dtos.CreateUserResponseDto, error) {
	emailExist, _ := service.Repository.FindByEmail(ctx, input.Email)
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
	fmt.Println("is active: ", input.IsActive)
	password, err := hash.HashingPassword(input.Password)
	if err != nil {
		return nil, err
	}
	record, err := service.Repository.Create(ctx, models.User{
		Password:    password,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
		Role:        input.Role,
		IsActive:    input.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return &dtos.CreateUserResponseDto{
		FirstName:   record.FirstName,
		LastName:    record.LastName,
		PhoneNumber: record.PhoneNumber,
		Email:       record.Email,
	}, nil
}
func (service serviceUser) GetAll(ctx context.Context, filter dtos.UserFilter) ([]models.User, *helper.PaginationMeta, error) {
	return service.Repository.GetAll(ctx, filter)
}

func (service serviceUser) FindById(ctx context.Context, id string) (*models.User, error) {
	return service.Repository.FindById(ctx, id)
}

func (service serviceUser) Update(ctx context.Context, input dtos.UpdateUserDto) error {
	err := service.Repository.Transaction(ctx, func(repo repositories.UserRepository) error {
		record, err := service.Repository.FindById(ctx, input.ID)
		if err != nil {
			return err
		}
		_ = service.Repository.Update(ctx,
			&models.User{
				ID: record.ID,
				BaseModel: shared.BaseModel{
					UpdatedDate: time.Now().Local(),
				},
				Password:      input.Password,
				FirstName:     input.FirstName,
				LastName:      input.LastName,
				PhoneNumber:   input.PhoneNumber,
				Email:         input.Email,
				IsActive:      input.IsActive,
				LastLoginDate: nil,
				LastLoginIP:   nil,
			},
		)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (service serviceUser) Delete(ctx context.Context, id string) error {
	user, _ := service.Repository.FindById(ctx, id)
	if user == nil {
		return &exception.ErrWithCode{
			Code: http.StatusNotFound,
			Err:  errors.New("user not found"),
		}
	}
	return service.Repository.Delete(ctx, id)
}

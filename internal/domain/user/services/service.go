package services

import (
	"context"
	"github.com/redis/go-redis/v9"
	"user-service/internal/domain/user/dtos"
	"user-service/internal/domain/user/models"
	"user-service/internal/domain/user/repositories"
	"user-service/pkg/helper"
)

type serviceUser struct {
	Repository repositories.UserRepository
	Redis      redis.Client
}

type UserService interface {
	// ----------private--------- \\

	Create(ctx context.Context, input dtos.CreateUserDto) (*dtos.CreateUserResponseDto, error)
	GetAll(ctx context.Context, filter dtos.UserFilter) ([]models.User, *helper.PaginationMeta, error)
	FindById(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, input dtos.UpdateUserDto) error
	Delete(ctx context.Context, id string) error

	// ----------public--------- \\

	Register(ctx context.Context, input dtos.RegisterDto) (*dtos.CreateUserResponseDto, error)
	Login(ctx context.Context, input dtos.LoginDto) (*dtos.TokenDto, error)
}

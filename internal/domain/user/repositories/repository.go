package repositories

import (
	"context"
	"gorm.io/gorm"
	"user-service/internal/domain/user/dtos"
	"user-service/internal/domain/user/models"
	"user-service/pkg/helper"
)

type repositoryUser struct {
	db *gorm.DB
}

type UserRepository interface {
	Transaction(ctx context.Context, fn func(repo UserRepository) error) error
	Create(ctx context.Context, user models.User) (*models.User, error)
	GetAll(ctx context.Context, filter dtos.UserFilter) ([]models.User, *helper.PaginationMeta, error)
	FindById(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phoneNumber string) (*models.User, error)
	Update(ctx context.Context, input *models.User) error
	Delete(ctx context.Context, id string) error
}

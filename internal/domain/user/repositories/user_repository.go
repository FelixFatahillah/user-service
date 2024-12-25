package repositories

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"user-service/internal/domain/user/dtos"
	"user-service/internal/domain/user/models"
	"user-service/pkg/exception"
	"user-service/pkg/helper"
	"user-service/pkg/logger"
)

func NewRepositoryUser(db *gorm.DB) *repositoryUser {
	return &repositoryUser{db}
}

func (repository repositoryUser) beginTransaction() *gorm.DB { return repository.db.Begin() }

func (repository repositoryUser) withTx(ctx context.Context, tx *gorm.DB) *repositoryUser {
	repository.db = tx.WithContext(ctx)
	return &repository
}

func (repository repositoryUser) Transaction(ctx context.Context, fn func(repo UserRepository) error) error {
	tx := repository.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	repo := repository.withTx(ctx, tx)
	err := fn(repo)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (repository repositoryUser) Create(ctx context.Context, user models.User) (*models.User, error) {
	if err := repository.db.WithContext(ctx).Create(&user).Error; err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return &user, nil
}

func (repository repositoryUser) GetAll(ctx context.Context, filter dtos.UserFilter) ([]models.User, *helper.PaginationMeta, error) {
	records := make([]models.User, 0)
	paginateMeta := helper.PaginationMeta{
		Page:  filter.Pagination.Page,
		Limit: filter.Pagination.Limit,
	}

	query := repository.db.WithContext(ctx).Model(models.User{}).Scopes(helper.PaginateScope(&filter.Pagination))
	query.Count(&paginateMeta.Total).Order("created_date DESC").Find(&records)

	paginateMeta.TotalPage = helper.GetTotalPage(paginateMeta.Total, paginateMeta.Limit)

	return records, &paginateMeta, nil
}

func (repository repositoryUser) FindById(ctx context.Context, id string) (*models.User, error) {
	record := &models.User{}
	err := repository.db.WithContext(ctx).Take(record, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &exception.ErrWithCode{
				Code: http.StatusBadRequest,
				Err:  errors.New(fmt.Sprintf("user %s not found", id)),
			}
		}
		return nil, err
	}

	return record, nil
}

func (repository repositoryUser) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	record := &models.User{}
	fmt.Println("email: ", email)
	err := repository.db.WithContext(ctx).Take(record, "email = ?", email).Error

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (repository repositoryUser) FindByPhone(ctx context.Context, phoneNumber string) (*models.User, error) {
	record := &models.User{}
	err := repository.db.WithContext(ctx).Take(record, "phone_number = ?", phoneNumber).Error

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (repository repositoryUser) Update(ctx context.Context, input *models.User) error {
	if err := repository.db.
		WithContext(ctx).
		Updates(input).Error; err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (repository repositoryUser) Delete(ctx context.Context, id string) error {
	err := repository.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&models.User{}).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

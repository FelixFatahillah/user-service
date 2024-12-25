package delivery

import (
	"gorm.io/gorm"
	repositoryUser "user-service/internal/domain/user/repositories"
)

type Repositories struct {
	UserRepository repositoryUser.UserRepository
}

func NewRepository(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepository: repositoryUser.NewRepositoryUser(db),
	}
}

package delivery

import (
	"github.com/redis/go-redis/v9"
	"user-service/internal/domain/user/services"
)

type Service struct {
	UserService services.UserService
}

type Deps struct {
	Repository *Repositories
	Redis      redis.Client
	//GRPC       *GRPC
}

func NewService(deps Deps) *Service {
	return &Service{
		UserService: services.NewServiceUser(deps.Repository.UserRepository, deps.Redis),
	}
}

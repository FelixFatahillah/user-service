package user

import (
	"github.com/gofiber/fiber/v2"
	"time"
	"user-service/internal/domain/user/services"
	"user-service/internal/middleware"
)

type handlerRESTUser struct {
	service services.UserService
}

func NewHandlerRESTUser(service services.UserService, router *fiber.App) {
	handler := handlerRESTUser{
		service,
	}

	routerV1 := router.Group("/api/v1")
	routerProtected := routerV1.Group("/private/users", middleware.RoleMiddleware())
	routerPublic := routerV1.Group("/public/users")

	routerProtected.Delete("/:id", handler.handlerDelete)
	routerProtected.Get("/:id", handler.handlerFindById)
	routerProtected.Put("/:id", handler.handlerUpdate)
	routerProtected.Post("/", handler.handlerCreate)
	routerProtected.Get("/", handler.handlerGetAll)

	routerPublic.Post("/register", middleware.RateLimit(30, 15*time.Second, "too many request"), handler.handlerRegister)
	routerPublic.Post("/login", middleware.RateLimit(30, 15*time.Second, "too many request"), handler.handlerLogin)

}

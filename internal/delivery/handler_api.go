package delivery

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"time"
	"user-service/internal/domain/user"
	"user-service/internal/domain/user/services"
	"user-service/pkg/exception"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type handler struct {
	serviceUser services.UserService
}

func NewHandler(serviceUser services.UserService) *handler {
	return &handler{
		serviceUser: serviceUser,
	}
}

const idleTimeout = 60 * time.Second

func (handler *handler) Init() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: exception.FiberErrorHandler,
		IdleTimeout:  idleTimeout,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(etag.New())
	app.Use(requestid.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	user.NewHandlerRESTUser(handler.serviceUser, app)

	return app
}

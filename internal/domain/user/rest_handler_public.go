package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"user-service/internal/domain/user/dtos"
	"user-service/pkg/response"
	"user-service/pkg/validation"
)

func (handler handlerRESTUser) handlerRegister(ctx *fiber.Ctx) error {
	var body dtos.RegisterDto
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}
	if err := validation.Validate(body); err != nil {
		return err
	}

	data, err := handler.service.Register(ctx.Context(), body)
	if err != nil {
		return err
	}

	return response.Respond(ctx, fiber.StatusOK, fmt.Sprintf("hello %s %s contact the administrator to actived your account", body.FirstName, *body.LastName), data, nil, nil)
}

func (handler handlerRESTUser) handlerLogin(ctx *fiber.Ctx) error {
	var body dtos.LoginDto
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}
	if err := validation.Validate(body); err != nil {
		return err
	}

	data, err := handler.service.Login(ctx.Context(), body)
	if err != nil {
		return err
	}

	return response.Respond(ctx, fiber.StatusOK, "success", data, nil, nil)
}

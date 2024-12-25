package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"user-service/internal/domain/user/dtos"
	"user-service/pkg/helper"
	"user-service/pkg/response"
	"user-service/pkg/validation"
)

func (handler handlerRESTUser) handlerCreate(ctx *fiber.Ctx) error {
	var body dtos.CreateUserDto
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}
	if err := validation.Validate(body); err != nil {
		return err
	}

	data, err := handler.service.Create(ctx.Context(), body)
	if err != nil {
		return err
	}

	return response.Respond(ctx, fiber.StatusOK, "success", data, nil, nil)
}

func (handler handlerRESTUser) handlerGetAll(ctx *fiber.Ctx) error {
	paginate := helper.Pagination{
		Page:  1,
		Limit: 10,
	}
	err := ctx.QueryParser(&paginate)
	if paginate.Limit >= 100 {
		paginate.Limit = 100
	}

	filter := dtos.UserFilter{
		Pagination: paginate,
	}
	_ = ctx.QueryParser(&filter)

	data, meta, err := handler.service.GetAll(ctx.Context(), filter)
	if err != nil {
		return err
	}

	return response.Respond(ctx, fiber.StatusOK, "success", data, nil, meta)
}

func (handler handlerRESTUser) handlerFindById(ctx *fiber.Ctx) error {
	data, err := handler.service.FindById(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}

	return response.Respond(ctx, fiber.StatusOK, "success", data, nil, nil)
}

func (handler handlerRESTUser) handlerUpdate(ctx *fiber.Ctx) error {
	var body dtos.UpdateUserDto
	body.ID = ctx.Params("id")
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}
	if err := validation.Validate(body); err != nil {
		return err
	}

	err := handler.service.Update(ctx.Context(), body)
	if err != nil {
		return err
	}

	return response.Respond(ctx, fiber.StatusOK, "success", nil, nil, nil)
}

func (handler handlerRESTUser) handlerDelete(ctx *fiber.Ctx) error {
	err := handler.service.Delete(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}

	return response.Respond(ctx, fiber.StatusOK, fmt.Sprintf("success deleted %s", ctx.Params("id")), nil, nil, nil)
}

package handler

import (
	"context"
	"errors"

	"github.com/Brigant/TestTask/app/config"
	"github.com/Brigant/TestTask/app/model"
	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
)

type AutService interface {
	GetToken(context.Context, model.IdentityData, config.AuthConfig) (model.Token, error)
}

type AuthHandler struct {
	Service    AutService
	AuthConfig config.AuthConfig
	Validator  *validator.Validate
}

func NewAuthHandler(service AutService, cfg config.AuthConfig, validator *validator.Validate) AuthHandler {
	return AuthHandler{
		Service:    service,
		AuthConfig: cfg,
		Validator:  validator,
	}
}

func (h AuthHandler) Login(ctx *fiber.Ctx) error {
	identity := model.IdentityData{}

	if err := ctx.BodyParser(&identity); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body parsing error")
	}

	if err := h.Validator.Struct(identity); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	token, err := h.Service.GetToken(ctx.UserContext(), identity, h.AuthConfig)
	if err != nil {
		if errors.Is(err, model.ErrTimeout) {
			return fiber.NewError(fiber.StatusRequestTimeout, err.Error())
		}

		if errors.Is(err, model.ErrNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(token)
}

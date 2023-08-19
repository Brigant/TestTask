package middleware

import (
	"strings"

	"github.com/Brigant/TestTask/app/config"
	"github.com/Brigant/TestTask/app/service"
	"github.com/gofiber/fiber/v2"
)

func NewUserIdentifier(cfg config.AuthConfig) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headers := ctx.GetReqHeaders()
		authorization := headers["Authorization"]
		if authorization == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "empty authorization header")
		}

		parts := strings.Split(authorization, " ")
		if parts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "authorization not bearer type")
		}

		userID, err := service.ParseToken(parts[1], cfg.SigningKey)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		ctx.Locals("userID", userID)

		return ctx.Next()
	}
}

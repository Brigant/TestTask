package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/Brigant/TestTask/app/api/handler"
	"github.com/Brigant/TestTask/app/api/middleware"
	"github.com/Brigant/TestTask/app/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

type Server struct {
	HTTPServer *fiber.App
}

func NewServer(cfg config.Config, handler handler.Handler) *Server {
	server := new(Server)

	fconfig := fiber.Config{
		ReadTimeout:  cfg.Server.AppReadTimeout,
		WriteTimeout: cfg.Server.AppWriteTimeout,
		IdleTimeout:  cfg.Server.AppIdleTimeout,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			ctx.Status(code)

			if err := ctx.JSON(err); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			return nil
		},
	}

	server.HTTPServer = fiber.New(fconfig)

	server.HTTPServer.Use(cors.New(corsConfig()))

	server.HTTPServer.Use(logger.New())

	server.HTTPServer.Use(recover.New())

	server.initRoutes(server.HTTPServer, handler, cfg)

	return server
}

func (s *Server) Shutdown(ctx context.Context) error {
	return fmt.Errorf("shutdown error: %w", s.HTTPServer.ShutdownWithContext(ctx))
}

func (s Server) initRoutes(app *fiber.App, h handler.Handler, cfg config.Config) {
	identifyUser := middleware.NewUserIdentifier(cfg.Auth)

	app.Static("/uploads", "./uploads")

	app.Post("/login", timeout.NewWithContext(h.Auth.Login, cfg.Server.AppReadTimeout))

	app.Get("/images", identifyUser, timeout.NewWithContext(h.Image.GetImages, cfg.Server.AppReadTimeout))
	app.Post("/upload-picture", identifyUser, timeout.NewWithContext(h.Image.SaveImage, cfg.Server.AppWriteTimeout))
}

func corsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     `*`,
		AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Credentials, Authorization",
		AllowMethods:     "GET, POST",
		AllowCredentials: true,
	}
}

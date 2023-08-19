package handler

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/Brigant/TestTask/app/model"
	"github.com/gofiber/fiber/v2"
)

type ImageService interface {
	GetImages(context.Context, string) ([]model.Image, error)
	SaveImage(context.Context, *multipart.Form, string) error
}

type ImageHandler struct {
	Service ImageService
}

// GetImages returns all images.
func NewImageHandler(service ImageService) ImageHandler {
	return ImageHandler{
		Service: service,
	}
}

// GetImages returns all images from storage
func (h ImageHandler) GetImages(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID")

	images, err := h.Service.GetImages(ctx.UserContext(), userID.(string))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(images)
}

// SaveImage saves a new image
func (h ImageHandler) SaveImage(ctx *fiber.Ctx) error {
	const (
		limitNumberItemsFile = 1
		imageFormKey         = "image"
	)
	userID := ctx.Locals("userID")

	allowedFileExtentions := []string{"jpg", "jpeg", "webp", "png"}

	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if len(form.File[imageFormKey]) < limitNumberItemsFile {
		return fiber.NewError(fiber.StatusBadRequest, "no attached image")
	}

	fileHeader := form.File[imageFormKey][0]

	if !isAllowedFileExtention(allowedFileExtentions, fileHeader.Filename) {
		return fiber.NewError(fiber.StatusBadRequest, "required format jpg/jpeg/webp/png")
	}

	if err := h.Service.SaveImage(ctx.UserContext(), form, userID.(string)); err != nil {
		if errors.Is(err, model.ErrDatabaseViolation) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON("success")
}

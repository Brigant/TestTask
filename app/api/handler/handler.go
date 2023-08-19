package handler

import (
	"strings"

	"github.com/Brigant/TestTask/app/config"
	"github.com/go-playground/validator/v10"
)

type ServiceInterfaces interface {
	ImageService
	AutService
}

type Handler struct {
	Image ImageHandler
	Auth  AuthHandler
}

func NewHandler(services ServiceInterfaces, cfg config.Config, validator *validator.Validate) Handler {
	return Handler{
		Image: NewImageHandler(services),
		Auth:  NewAuthHandler(services, cfg.Auth, validator),
	}
}

func isAllowedFileExtention(allowedList []string, fileName string) bool {
	nameParts := strings.Split(fileName, ".")

	fileExt := nameParts[len(nameParts)-1]
	for _, i := range allowedList {
		if i == fileExt {
			return true
		}
	}

	return false
}

func symbolsCounter(sentence string) int {
	runes := []rune(sentence)

	return len(runes)
}

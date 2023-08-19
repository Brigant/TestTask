package service

import "github.com/Brigant/TestTask/app/config"

const (
	uploadDirectory = "uploads/"
	filePermition   = 0o666
)

type Storager interface {
	ImageStorager
	UserStorager
}

type Services struct {
	ImageService
	AuthService
}

func New(storage Storager, cfg config.Server) Services {
	return Services{
		ImageService{Storage: storage, cfg: cfg},
		AuthService{Storage: storage},
	}
}

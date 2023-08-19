package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Brigant/TestTask/app/config"
	"github.com/Brigant/TestTask/app/model"
)

type ImageStorager interface {
	InsertImage(context.Context, model.Image) error
	SelectImages(ctx context.Context, userID string) ([]model.Image, error)
}

type ImageService struct {
	Storage ImageStorager
	cfg     config.Server
}

func (s ImageService) GetImages(ctx context.Context, userID string) ([]model.Image, error) {
	images, err := s.Storage.SelectImages(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service GetImages error: %w", err)
	}

	return images, nil
}

func (s ImageService) SaveImage(ctx context.Context, form *multipart.Form, userID string) error {
	fileHeader := form.File["image"][0]

	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("service SaveImage error: %w", err)
	}

	baseUrl := "http://" + s.cfg.AppHost + ":" + s.cfg.AppPort

	image := model.Image{
		UserID:    userID,
		ImagePath: uploadDirectory + makeUniqueFileName(fileHeader.Filename),
		ImageURL:  baseUrl + "/" + uploadDirectory + makeUniqueFileName(fileHeader.Filename),
	}

	osFile, err := os.OpenFile(image.ImagePath, os.O_WRONLY|os.O_CREATE, filePermition)
	if err != nil {
		return fmt.Errorf("service SaveImage error: %w", err)
	}

	defer osFile.Close()

	written, err := io.Copy(osFile, file)
	if err != nil {
		return fmt.Errorf(" written bytes: %v, service SaveImage error: %w", written, err)
	}

	if err := s.Storage.InsertImage(ctx, image); err != nil {
		if err := os.Remove(image.ImagePath); err != nil {
			return fmt.Errorf("service SaveImage error: %w", err)
		}

		return fmt.Errorf("service SaveImage error: %w", err)
	}

	return nil
}

func makeUniqueFileName(fileName string) string {
	nameParts := strings.Split(fileName, ".")

	fileExt := nameParts[len(nameParts)-1]

	unixtime := time.Now().Unix()
	strUnixTime := strconv.Itoa(int(unixtime))

	uniqueFileName := strings.TrimRight(fileName, "."+fileExt) + "_" + strUnixTime + "." + fileExt

	return uniqueFileName
}

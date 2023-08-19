package storage

import (
	"context"

	"github.com/Brigant/TestTask/app/model"
)

func (s Storage) InsertImage(ctx context.Context, image model.Image) error {
	query := `INSERT INTO public.images ("user_id", "image_path", "image_url")
			VALUES (:user_id, :image_path, :image_url)`

	sqlResult, err := s.db.NamedExecContext(ctx, query, image)

	return handleExexError("storage InsertImage error", sqlResult, err)
}

func (s Storage) SelectImages(ctx context.Context, userID string) ([]model.Image, error) {
	query := `SELECT "id", "user_id", "image_path", "image_url"
				FROM public.images WHERE "user_id"=$1;
			`
	images := []model.Image{}

	if err := s.db.SelectContext(ctx, &images, query, userID); err != nil {
		return nil, handleSelectError("storage SelectImagea error", err)
	}

	return images, nil
}

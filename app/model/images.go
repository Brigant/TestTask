package model

type Image struct {
	ID        string `json:"id" db:"id" validate:"omitempty,uuid"`
	UserID    string `json:"user_id" db:"user_id" validate:"omitempty,uuid"`
	ImagePath string `json:"image_path" db:"image_path" validate:"required"`
	ImageURL  string `json:"image_url" db:"image_url" validate:"required"`
}

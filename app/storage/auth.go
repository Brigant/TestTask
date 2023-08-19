package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Brigant/TestTask/app/model"
	"github.com/lib/pq"
)

func (r Storage) FindUsernameWithPasword(ctx context.Context, userName, password string) (string, error) {
	var userID string

	query := `SELECT id FROM users WHERE username=$1 and password_hash=$2`

	if err := r.db.GetContext(ctx, &userID, query, userName, password); err != nil {
		pqErr := new(pq.Error)
		if errors.As(err, &pqErr) && pqErr.Code.Name() == "query_canceled" {
			return "", model.ErrTimeout
		}

		if errors.Is(err, sql.ErrNoRows) {
			return "", model.ErrNotFound
		}

		return "", fmt.Errorf("error while GetContext(): %w", err)
	}

	return userID, nil
}

func SHA256(password, salt string) string {
	sum := sha256.Sum256([]byte(password + salt))

	return fmt.Sprintf("%x", sum)
}

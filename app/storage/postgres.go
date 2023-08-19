package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Brigant/TestTask/app/config"
	"github.com/Brigant/TestTask/app/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	minEexpectedAffectedRows = 1
	codeUniqueViolation      = "unique_violation"
	codeNoData               = "no_data"
	codeForeignKeyViolation  = "foreign_key_violation"
	codeCanceled             = "query_canceled"
)

type Storage struct {
	db *sqlx.DB
}

// Returns an entety of the Storage.
func New(cfg config.Config) (Storage, error) {
	database, err := sqlx.Connect("postgres", fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Database, cfg.DB.Password, cfg.DB.SSLmode))
	if err != nil {
		return Storage{}, fmt.Errorf("cannot connect to db: %w", err)
	}

	return Storage{db: database}, nil
}

// Close closes the database connection.
func (s Storage) Close() error {
	return fmt.Errorf("error hapens while db.close: %w", s.db.Close())
}

func handleExexError(message string, res sql.Result, err error) error {
	if err != nil {
		pqError := new(pq.Error)

		if errors.As(err, &pqError) {
			switch pqError.Code.Name() {
			case codeCanceled:
				return model.ErrTimeout
			case codeUniqueViolation:
				return fmt.Errorf("%s: %w", message, model.ErrDatabaseViolation)
			case codeForeignKeyViolation:
				return fmt.Errorf("%s: %w", message, model.ErrDatabaseViolation)
			default:
				return fmt.Errorf("%s: %w", message, err)
			}
		}
	}

	number, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", message, model.ErrDatabaseViolation)
	}

	if number < minEexpectedAffectedRows {
		return fmt.Errorf("%s: %w", message, model.ErrNothinChanged)
	}

	return nil
}

func handleSelectError(message string, err error) error {
	pqError := new(pq.Error)

	if errors.As(err, &pqError) && pqError.Code.Name() == codeCanceled {
		return model.ErrTimeout
	}

	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}

	return fmt.Errorf("%s: %w", message, err)
}

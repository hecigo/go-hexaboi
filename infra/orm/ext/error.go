package ext

import (
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

// Used to format any error to HPIResult
func Errorf(err error) (error, bool) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &core.HPIResult{
			Status:    404,
			Message:   err.Error(),
			ErrorCode: "ORM_ERROR",
		}, true
	}

	if errors.Is(err, gorm.ErrInvalidField) || errors.Is(err, gorm.ErrEmptySlice) {
		return &core.HPIResult{
			Status:    400,
			Message:   err.Error(),
			ErrorCode: "ORM_ERROR",
		}, true
	}

	// GORM doesn't override dialect errors, have to detect manually.
	if err, ok := err.(*pgconn.PgError); ok {
		return &core.HPIResult{
			Status:    500,
			Message:   err.Error(),
			ErrorCode: "ORM_ERROR",
		}, true
	}

	return err, false
}

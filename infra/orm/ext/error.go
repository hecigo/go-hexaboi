package ext

import (
	"errors"

	"github.com/hecigo/goutils"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// Used to format any error to HPIResult
func Errorf(err error) (error, bool) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &goutils.APIRes{
			Status:    404,
			Message:   err.Error(),
			ErrorCode: "ORM_ERROR",
		}, true
	}

	if errors.Is(err, gorm.ErrInvalidField) || errors.Is(err, gorm.ErrEmptySlice) {
		return &goutils.APIRes{
			Status:    400,
			Message:   err.Error(),
			ErrorCode: "ORM_ERROR",
		}, true
	}

	// GORM doesn't override dialect errors, have to detect manually.
	if err, ok := err.(*pgconn.PgError); ok {
		return &goutils.APIRes{
			Status:    500,
			Message:   err.Error(),
			ErrorCode: "ORM_ERROR",
		}, true
	}

	if err, ok := err.(*goutils.APIRes); ok {
		return err, true
	}

	return err, false
}

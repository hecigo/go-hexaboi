package orientdb

import (
	"strconv"

	"github.com/hecigo/goutils"
)

type Errors struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Reason  int    `json:"reason"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

func (e Error) Error() string {
	return e.Content
}

func (e Error) ToHPIError() error {
	return &goutils.APIRes{
		Status:    e.Code,
		Message:   strconv.Itoa(e.Reason),
		ErrorCode: "ORM_ERROR",
		Data:      e.Content,
	}
}

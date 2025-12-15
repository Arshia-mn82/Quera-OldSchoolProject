package router

import (
	"OldSchool/internal/service"
	"OldSchool/internal/transport/protocol"
	"errors"
)

func fromServiceError(err error) protocol.Response {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return protocol.Response{Status: false, Message: "invalid input", Data: nil}
	case errors.Is(err, service.ErrNotFound):
		return protocol.Response{Status: false, Message: "not found", Data: nil}
	case errors.Is(err, service.ErrRoleMismatch):
		return protocol.Response{Status: false, Message: "role mismatch", Data: nil}
	case errors.Is(err, service.ErrDuplicateEnrollment):
		return protocol.Response{Status: false, Message: "duplicate enrollment", Data: nil}
	case errors.Is(err, service.ErrDifferentSchool):
		return protocol.Response{Status: false, Message: "different school not allowed", Data: nil}
	case errors.Is(err, service.ErrSchoolAlreadyExists):
		return protocol.Response{Status: false, Message: "school already exists", Data: nil}
	default:
		return protocol.Response{Status: false, Message: "internal error", Data: nil}
	}
}

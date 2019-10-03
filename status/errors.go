// Package errors provides a way to return detailed information
// for an RPC request error. The error is normally JSON encoded.
package status

import (
	"fmt"
	proto "github.com/mhconradt/proto/status"
)

var SuccessValue = &proto.Status{
	Code: 200,
	Message: "success",
}

// Error implements the error interface.
// New generates a custom error.
func New(id, detail string, code int32) *proto.Status {
	return &proto.Status{
		Code:   code,
		Message: detail,
	}
}

func Success() *proto.Status {
	return &proto.Status{
		Code: 200,
		Message: "success",
	}
}

// BadRequest generates a 400 error.
func BadRequest(format string, a ...interface{}) *proto.Status {
	return &proto.Status{
		Code:   int32(400),
		Message: fmt.Sprintf(format, a...),
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(format string, a ...interface{}) *proto.Status {
	return &proto.Status{
		Code:   401,
		Message: fmt.Sprintf(format, a...),
	}
}

// Forbidden generates a 403 error.
func Forbidden(format string, a ...interface{}) *proto.Status {
	return &proto.Status{
		Code:   403,
		Message: fmt.Sprintf(format, a...),
	}
}

// NotFound generates a 404 error.
func NotFound(format string, a ...interface{}) *proto.Status {
	return &proto.Status{
		Code:   404,
		Message: fmt.Sprintf(format, a...),
	}
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(format string, a ...interface{}) *proto.Status {
	return &proto.Status{
		Code:   405,
		Message: fmt.Sprintf(format, a...),
	}
}

// Timeout generates a 408 error.
func Timeout(format string, a ...interface{}) *proto.Status {
	return &proto.Status{
		Code:   408,
		Message: fmt.Sprintf(format, a...),
	}
}

// Conflict generates a 409 error.
func Conflict(format string, a ...interface{}) *proto.Status {
	return &proto.Status{
		Code:   409,
		Message: fmt.Sprintf(format, a...),
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(format string, a ...interface{}) *proto.Status {
	return &proto.Status{
		Code:   500,
		Message: fmt.Sprintf(format, a...),
	}
}

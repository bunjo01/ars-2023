package tracer

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
)

// swagger:response ErrorResponse
type ErrorResponse struct {
	// Error status code
	// in: int64
	Status int `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("error code: %d\n%s", e.Status, e.Message)
}

var codes = map[int]string{
	400: "bad request",
	401: "unauthorized",
	403: "forbidden",
	404: "not found",
	405: "not allowed",
	406: "not accepted",
	409: "conflict",
	415: "unsupported media type",
	418: "I'm a teapot",
}

func NewError(code int, span opentracing.Span) *ErrorResponse {
	er := &ErrorResponse{
		Status:  code,
		Message: codes[code],
	}
	LogError(span, *er)
	return er
}

package errors

// swagger:response ErrorResponse
type ErrorResponse struct {
	// Error status code
	// in: int64
	Status int `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
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

func NewError(code int) *ErrorResponse {
	return &ErrorResponse{
		Status:  code,
		Message: codes[code],
	}
}

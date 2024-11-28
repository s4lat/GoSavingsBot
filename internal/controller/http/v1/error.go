package v1

import "fmt"

type ErrorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(errMsg string) ErrorResponse {
	return ErrorResponse{
		Error: errMsg,
	}
}

func newErrorResponseF(errMsg string, args ...any) ErrorResponse {
	return ErrorResponse{
		Error: fmt.Sprintf(errMsg, args...),
	}
}

var (
	internalError = ErrorResponse{Error: "something went wrong"}
)

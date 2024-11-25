package v1

type ErrorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(errMsg string) ErrorResponse {
	return ErrorResponse{
		Error: errMsg,
	}
}

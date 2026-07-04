package web

// ErrorDto is a standard HTTP error response body.
type ErrorDto struct {
	Error string `json:"error"`
}

// NewErrorDto wraps an error into an ErrorDto.
func NewErrorDto(err error) ErrorDto {
	return ErrorDto{Error: err.Error()}
}

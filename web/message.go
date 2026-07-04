package web

// MessageDto is a standard HTTP message response body.
type MessageDto struct {
	Message string `json:"message"`
}

// NewMessageDto wraps a string message into a MessageDto.
func NewMessageDto(msg string) MessageDto {
	return MessageDto{Message: msg}
}

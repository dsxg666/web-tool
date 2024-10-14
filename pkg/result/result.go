package result

type Result[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func NewResult[T any](code, message string, data T) *Result[T] {
	return &Result[T]{Code: code, Message: message, Data: data}
}

package result

var (
	Success            = NewResult("100", "Success", "")
	InvalidRequestData = NewResult("200", "Invalid request data!", "")
	NoJwtToken         = NewResult("201", "Authorization header is missing", "")
	InvalidJwtToken    = NewResult("202", "Invalid jwt token", "")
	FileNotFound       = NewResult("203", "File not found", "")
	UserIdNotFound     = NewResult("204", "UserId not found", "")

	InternalServerError = NewResult("500", "Internal server error", "")
)

func SuccessWithData[T any](data T) *Result[T] {
	return NewResult("100", "Success", data)
}

func SuccessWithMessage[T any](info string, data T) *Result[T] {
	return NewResult("100", info, data)
}

func OperateError[T any](cause string, data T) *Result[T] {
	return NewResult("300", cause, data)
}

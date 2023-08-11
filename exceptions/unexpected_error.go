package exceptions

const (
	UnexpectedErrorMessage = "Ocorreu um erro inesperado."
)

type UnexpectedError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewUnexpectedError(message string) UnexpectedError {
	return UnexpectedError{
		Message: message,
		Code:    500,
	}
}

func (b UnexpectedError) Error() string {
	return b.Message
}

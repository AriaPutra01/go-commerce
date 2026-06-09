package exception

type AppError struct {
	StatusCode int    `json:"-"`
	ErrCode    string `json:"code"`
	Message    string `json:"message"`
	Errors     string `json:"errors,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

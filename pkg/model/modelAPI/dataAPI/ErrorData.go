package dataAPI

type ErrorData struct {
	Error string `json:"error"`
}

// NewErrorData - constructor of ErrorData struct.
func NewErrorData(err string) *ErrorData {
	return &ErrorData{
		Error: err,
	}
}

// GetContent - return the data of current context.
func (data *ErrorData) GetContent() interface{} {
	return data.Error
}

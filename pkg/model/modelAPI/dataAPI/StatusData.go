package dataAPI

const SuccessStatus = "success"
const FailureStatus = "failure"

type StatusData struct {
	Status string `json:"status"`
}

// NewStatusData - constructor of StatusData struct.
func NewStatusData(isSuccess bool) *StatusData {
	var status string

	if isSuccess {
		status = SuccessStatus
	} else {
		status = FailureStatus
	}

	return &StatusData{
		Status: status,
	}
}

// GetContent - return the data of current context.
func (data *StatusData) GetContent() interface{} {
	return data.Status
}

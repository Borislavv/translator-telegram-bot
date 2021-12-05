package model

type UpdatedMessages struct {
	Messages []UpdatedMessage `json:"result"`
}

type UpdatedMessage struct {
	QueueId int64 `json:"update_id"`
	Data    struct {
		Chat struct {
			ID       int64  `json:"id"`
			Title    string `json:"first_name"`
			Username string `json:"username"`
		} `json:"chat"`
		Date int64  `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

// NewUpdatedMessages - UpdatedMessages struct creator
func NewUpdatedMessages() *UpdatedMessages {
	return &UpdatedMessages{}
}

// NewUpdatedMessages - UpdatedMessage struct creator
func NewUpdatedMessage() *UpdatedMessage {
	return &UpdatedMessage{}
}

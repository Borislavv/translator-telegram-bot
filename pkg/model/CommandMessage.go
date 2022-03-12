package model

type CommandMessage struct {
	OriginMessage *UpdatedMessage
}

// NewCommandMessage - constructor of CommandMessage structure.
func NewCommandMessage(originMessage *UpdatedMessage) *CommandMessage {
	return &CommandMessage{
		OriginMessage: originMessage,
	}
}

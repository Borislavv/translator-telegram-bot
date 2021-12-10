package model

type TelegramResponseMessage struct {
	ChatId  string
	Message string
}

// NewTelegramResponseMessage - constructor of TelegramResponseMessage
func NewTelegramResponseMessage(chatId string, message string) *TelegramResponseMessage {
	return &TelegramResponseMessage{
		ChatId:  chatId,
		Message: message,
	}
}

// GetChatId - getter of prop. ChatId
func (message *TelegramResponseMessage) GetChatId() string {
	return message.ChatId
}

// GetMessage - getter of prop. Message
func (message *TelegramResponseMessage) GetMessage() string {
	return message.Message
}

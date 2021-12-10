package model

type TelegramResponseMessage struct {
	chatId string
	text   string
}

// NewTelegramResponseMessage - constructor of TelegramResponseMessage
func NewTelegramResponseMessage(chatId string, message string) *TelegramResponseMessage {
	return &TelegramResponseMessage{
		chatId: chatId,
		text:   message,
	}
}

// GetChatId - getter of prop. ChatId
func (message *TelegramResponseMessage) GetChatId() string {
	return message.chatId
}

// GetMessage - getter of prop. Message
func (message *TelegramResponseMessage) GetMessage() string {
	return message.text
}

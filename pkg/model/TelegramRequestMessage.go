package model

type TelegramRequestMessage struct {
	offset int64
}

// NewTelegramRequestMessage - constructor of TelegramRequestMessage
func NewTelegramRequestMessage(offset int64) *TelegramRequestMessage {
	return &TelegramRequestMessage{
		offset: offset,
	}
}

// GetOffset - getter of prop. ChatId
func (message *TelegramRequestMessage) GetOffset() int64 {
	return message.offset
}

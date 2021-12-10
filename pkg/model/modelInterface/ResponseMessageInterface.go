package modelInterface

type ResponseMessageInterface interface {
	GetChatId() string
	GetMessage() string
}

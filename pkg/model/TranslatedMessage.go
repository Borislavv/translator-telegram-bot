package model

type TranslatedMessage struct {
	Translations []Translation `json:"sentences"`
}

type Translation struct {
	Text string `json:"trans"`
}

// NewTranslatedMessage - TranslatedMessage struct creator
func NewTranslatedMessage() *TranslatedMessage {
	return &TranslatedMessage{}
}

package translator

import (
	"context"

	"github.com/Borislavv/Translator-telegram-bot/pkg/api/grpc/service/translatorGRPCInterface"
)

type TranslatorGRPC struct {
	translator *TranslatorService
}

// NewTranslatorGRPC - constructor of TranslatorGRPC structure.
func NewTranslatorGRPC(translator *TranslatorService) *TranslatorGRPC {
	return &TranslatorGRPC{
		translator: translator,
	}
}

// Translate - translate target text to opposite languagte (en/ru).
func (server *TranslatorGRPC) Translate(
	ctx context.Context,
	req *translatorGRPCInterface.RequestForTranslation,
) (*translatorGRPCInterface.TranslationResponse, error) {
	translation, err := server.translator.TranslateText(req.GetText())
	if err != nil {
		return nil, err
	}

	return &translatorGRPCInterface.TranslationResponse{
		Translation: translation,
	}, nil
}

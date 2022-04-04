package translator

import (
	"strings"
	"testing"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
)

func TestTranslator(t *testing.T) {
	translator := NewTranslatorService(
		NewTranslatorGateway(
			manager.New(
				config.New().Load(),
			),
		),
	)

	// test of `TranslateText` method
	t.Run("TranslateText", func(t *testing.T) {
		expected := "привет"

		translation, err := translator.TranslateText("hello")
		if err != nil {
			t.Error(err)
		}

		// excape case sensetivity
		translation = strings.ToLower(translation)

		if translation != expected {
			t.Errorf("expected value %v doesn't match the recived %v", expected, translation)
		}
	})

	// test of `detectLanguages` method
	t.Run("detectLanguages", func(t *testing.T) {
		ruText := "привет"

		detectedRuLanguage, _ := translator.detectLanguages(ruText)
		if detectedRuLanguage != RuLanguage {
			t.Errorf("expected value %v doesn't match the recived %v", RuLanguage, detectedRuLanguage)
		}
	})
}

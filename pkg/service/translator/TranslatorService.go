package translator

import (
	"log"
	"regexp"

	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

// DetectLanguages - detecting target and source language by regex
func DetectLanguages(text string) (string, string) {
	matched, err := regexp.MatchString("~[а-яА-Я]+~", text)
	if err != nil {
		log.Fatalln(util.Trace() + err.Error())
		return EnLanguage, RuLanguage
	}

	if matched {
		return RuLanguage, EnLanguage
	}

	return EnLanguage, RuLanguage
}

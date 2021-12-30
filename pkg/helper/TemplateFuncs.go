package helper

import "regexp"

type TemplateFuncsHelper struct {
}

// NewTemplateFuncsHelper - constructor of NewTemplateFuncsHelper struct
func NewTemplateFuncsHelper() *TemplateFuncsHelper {
	return &TemplateFuncsHelper{}
}

// Inc - increment function which will be called into the tempalates
func (funcs *TemplateFuncsHelper) Inc(i int) int {
	return i + 1
}

// PregReplace - simple preg_replace function
func (funcs *TemplateFuncsHelper) PregReplace(pattern string, subject string, replacement string) string {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return subject
	}

	return reg.ReplaceAllString(subject, replacement)
}

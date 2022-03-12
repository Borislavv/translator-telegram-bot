package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ExctractId - get the first pattern:`\/int`. value from target string.
func ExctractId(str string) (int64, error) {
	// receiveing id of target notification
	rgx, err := regexp.Compile(`\/\d+`)
	if err != nil {
		return 0, err
	}

	foundStr := rgx.FindString(str)

	var id int
	if foundStr != "" {
		id, err = strconv.Atoi(strings.Trim(foundStr, "/"))
		if err != nil {
			return 0, err
		} else if id == 0 {
			return 0, errors.New("Cannot exctract identifier from target string.")
		}
	}

	return int64(id), nil
}

// ExtractInt64 - get the first int. value from target string.
func ExtractInt64(str string) (int64, error) {
	// receiveing id of target notification
	rgx, err := regexp.Compile(`\d+`)
	if err != nil {
		return 0, err
	}

	foundStr := rgx.FindString(str)

	var intVal int
	if foundStr != "" {
		intVal, err = strconv.Atoi(foundStr)
		if err != nil {
			return 0, err
		} else if intVal == 0 {
			return 0, errors.New("Cannot exctract int64 from target string (zero(0) is unavailable).")
		}
	}

	return int64(intVal), nil
}

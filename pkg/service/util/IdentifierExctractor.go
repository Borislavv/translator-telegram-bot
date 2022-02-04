package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ExctractId - get the first int. value from target string.
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

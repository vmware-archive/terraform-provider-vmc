package vmc

import (
	"strings"
)

func IsValidString(str string) (b bool) {

	if len(strings.TrimSpace(str)) == 0 {
		return false
	}
	return true
}
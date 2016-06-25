package wain

import (
	"regexp"
	"strconv"
)

func StringToInt(value string) int {
	reg, _ := regexp.Compile("[^0-9]+")
	value = reg.ReplaceAllString(value, "")
	number, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0
	}
	return int(number)
}

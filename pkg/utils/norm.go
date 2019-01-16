package utils

import (
	"strconv"
	"strings"
)

func NumberGroupStringNormalize(input string) (string, []uint) {
	if input == "0" {
		return "", []uint{}
	}
	strSlice := strings.Split(input, ",")

	var resultSlice []string
	var idSlice []uint

	for _, strNum := range strSlice {
		if id, err := strconv.Atoi(strNum); err == nil {
			resultSlice = append(resultSlice, strNum)
			idSlice = append(idSlice, uint(id))
		}
	}
	return strings.Join(resultSlice, ","), idSlice
}

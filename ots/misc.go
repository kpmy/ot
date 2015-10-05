package ots

import (
	"strconv"
	"unicode"
)

func Token2(r rune) string {
	return string([]rune{r})
}

func Token(r rune) string {
	if unicode.IsSpace(r) || int(r) <= int(' ') {
		return strconv.Itoa(int(r)) + "U"
	} else {
		return string([]rune{r})
	}
}

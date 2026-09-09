package mcpobserve

import (
	"strings"
	"unicode/utf8"
)

func boundedText(value string, maximum int) (string, bool) {
	value = strings.ToValidUTF8(value, "\uFFFD")
	if len(value) <= maximum {
		return value, false
	}
	cut := maximum
	for cut > 0 && !utf8.RuneStart(value[cut]) {
		cut--
	}
	return value[:cut], true
}

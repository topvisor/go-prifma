package utils

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func UnescapeString(s string) string {
	if !strings.ContainsRune(s, '\\') && utf8.ValidString(s) {
		return s
	}

	buf := new(strings.Builder)
	buf.Grow(3 * len(s) / 2) // Try to avoid more allocations.

	var c rune
	var tail string
	var err error

	for len(s) > 0 {
		if c, _, tail, err = strconv.UnquoteChar(s, 0); err != nil {
			buf.WriteByte(s[0])
			s = s[1:]
		} else {
			buf.WriteRune(c)
			s = tail
		}
	}

	return buf.String()
}

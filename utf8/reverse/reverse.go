//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))

	cursor := 0
	for i := utf8.RuneCountInString(input); i > 0; i-- {
		r, size := utf8.DecodeLastRuneInString(input[:len(input)-cursor])
		builder.WriteRune(r)
		cursor += size
	}

	return builder.String()
}

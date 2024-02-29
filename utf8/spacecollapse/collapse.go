//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))

	cursor := 0
	prevSpace := false

	for i := utf8.RuneCountInString(input); i > 0; i-- {
		r, size := utf8.DecodeRuneInString(input[cursor:])
		cursor += size

		isSpace := unicode.IsSpace(r)
		if !isSpace {
			builder.WriteRune(r)
		} else if !prevSpace {
			builder.WriteRune(' ')
		}

		prevSpace = isSpace
	}

	return builder.String()
}

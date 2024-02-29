//go:build !solution

package varfmt

import (
	"container/list"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Stack struct {
	internal *list.List
}

func (c *Stack) Push(v any) {
	c.internal.PushFront(v)
}

func (c *Stack) Pop() {
	if c.internal.Len() == 0 {
		panic("stack is empty")
	}

	c.internal.Remove(c.internal.Front())
}

func (c *Stack) Front() *list.Element {
	return c.internal.Front()
}

func (c *Stack) Len() int {
	return c.internal.Len()
}

func countBalancedBrackets(s string) int {
	stack := &Stack{internal: list.New()}
	count := 0

	for _, ch := range s {
		switch ch {
		case '{':
			stack.Push(ch)
		case '}':
			if stack.Len() > 0 {
				front := stack.Front()
				if front != nil && front.Value.(rune) == '{' {
					stack.Pop()
					count++
				}
			}
		}
	}

	return count
}

func Sprintf(format string, args ...interface{}) string {
	// brackets := countBalancedBrackets(format)
	brackets := strings.Count(format, "{") // Assume balanced to pass perf test

	if brackets == 0 {
		return format
	}

	argsStrings := make([]string, 0, len(args))
	argsMaxLen := 0

	for _, arg := range args {
		s := fmt.Sprint(arg)
		argsStrings = append(argsStrings, s)
		argsMaxLen = max(len(s), argsMaxLen)
	}

	var builder strings.Builder
	builder.Grow(len(format) - brackets*2 + argsMaxLen*brackets)

	pos := 0
	cursor := 0
	number := 0
	numberLen := -1

	for i := utf8.RuneCountInString(format); i > 0; i-- {
		r, size := utf8.DecodeRuneInString(format[cursor:])
		cursor += size

		switch true {
		case r == '{':
			numberLen = 0
		case r == '}':
			index := number
			if numberLen <= 0 {
				index = pos
			}

			builder.WriteString(argsStrings[index])

			number = 0
			numberLen = -1
			pos += 1
		case numberLen >= 0:
			number = number*10 + (int(r) - '0')
			numberLen += 1
		default:
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

//go:build !solution

package speller

import "math"

func pow(p int) int64 {
	return int64(math.Pow(1000, float64(p)))
}

func Spell(n int64) string {
	if n < 0 {
		return "minus " + Spell(n*-1)
	}

	if n == 0 {
		return "zero"
	}

	till19 := []string{
		"one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen",
		"seventeen", "eighteen", "nineteen",
	}

	if n < 20 {
		return till19[n-1]
	}

	tens := []string{"twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}

	spellSeparated := func(num int64, delim string) string {
		if num == 0 {
			return ""
		}

		return delim + Spell(num)
	}

	if n < 100 {
		return tens[n/10-2] + spellSeparated(n%10, "-")
	}

	if n < 1000 {
		return till19[n/100-1] + " hundred" + spellSeparated(n%100, " ")
	}

	for i, w := range []string{"thousand", "million", "billion"} {
		p := i + 1

		if n < pow(p+1) {
			return Spell(n/pow(p)) + " " + w + spellSeparated(n%pow(p), " ")
		}
	}

	panic("unreachable")
}

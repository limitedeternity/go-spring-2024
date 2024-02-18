//go:build !solution

package main

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	files := os.Args[1:]

	if len(files) == 0 {
		return
	}

	result := make(map[string]int64)

	for _, file := range files {
		data, err := os.ReadFile(file)
		check(err)

		lines := strings.Split(string(data), "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)

			elem := result[line]
			result[line] = elem + 1
		}
	}

	for key, value := range result {
		if value == 1 {
			continue
		}

		fmt.Printf("%d\t%s\n", value, key)
	}
}

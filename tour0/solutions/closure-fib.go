//go:build exclude

package solutions

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	storage := [2]int{0, 1}
	index := 0

	return func() int {
		prevIndex := index

		index = (index + 1) % len(storage)
		storage[index] = sum(storage[:])

		return storage[prevIndex]
	}
}

func sum(numbers []int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

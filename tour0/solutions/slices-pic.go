//go:build exclude

package solutions

import (
	"fmt"

	"golang.org/x/exp/constraints"
	"golang.org/x/tour/pic"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func Pic(dx, dy int) [][]uint8 {
	ys := make([][]uint8, dy)

	for i := 0; i < dy; i++ {
		ysI := make([]uint8, dx)

		for j := 0; j < dx; j++ {
			ysI[j] = uint8(i * j)
		}

		ys[i] = ysI
	}

	return ys
}

func printSlice[T Number](s []T) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func main() {
	pic.Show(Pic)
}

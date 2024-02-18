//go:build exclude

package solutions

import (
	"fmt"
	"math"
)

func Sqrt(n float64) float64 {
	var root float64

	for x := n; ; {
		root = 0.5 * (x + (n / x))

		if math.Abs(root-x) < 0.000001 {
			break
		}

		x = root
	}

	return root
}

func main() {
	fmt.Println(Sqrt(2))
}

//go:build exclude

package solutions

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(n float64) (float64, error) {
	var root float64

	if n < 0 {
		return root, ErrNegativeSqrt(n)
	}

	for x := n; ; {
		root = 0.5 * (x + (n / x))

		if math.Abs(root-x) < 0.000001 {
			break
		}

		x = root
	}

	return root, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}

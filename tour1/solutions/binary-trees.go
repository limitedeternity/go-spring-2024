//go:build exclude

package solutions

import (
	"fmt"

	"golang.org/x/tour/tree"
)

func walkImpl(t *tree.Tree, ch chan int) {
	if t.Left != nil {
		walkImpl(t.Left, ch)
	}

	ch <- t.Value

	if t.Right != nil {
		walkImpl(t.Right, ch)
	}
}

func Walk(t *tree.Tree, ch chan int) {
	walkImpl(t, ch)
	close(ch)
}

func Same(t1, t2 *tree.Tree) bool {
	c1, c2 := make(chan int), make(chan int)

	go Walk(t1, c1)
	go Walk(t2, c2)

	for {
		v1, ok1 := <-c1
		v2, ok2 := <-c2

		if v1 != v2 || ok1 != ok2 {
			return false
		}

		if !ok1 && !ok2 {
			return true
		}
	}
}

func main() {
	fmt.Print("Same(tree.New(1), tree.New(1)): ")

	if Same(tree.New(1), tree.New(1)) {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}

	fmt.Print("Same(tree.New(1), tree.New(2)): ")

	if Same(tree.New(1), tree.New(2)) {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}

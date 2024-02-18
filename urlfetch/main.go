//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	urls := os.Args[1:]

	if len(urls) == 0 {
		return
	}

	for _, url := range urls {
		func() {
			resp, err := http.Get(url)
			check(err)

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				check(err)
			}(resp.Body)

			body, err := io.ReadAll(resp.Body)
			check(err)

			fmt.Println(string(body))
		}()
	}
}

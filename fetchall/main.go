//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func protect(msgChan chan string) func(func(chan string, string), string) {
	return func(f func(chan string, string), url string) {
		defer func() {
			if err := recover(); err != nil {
				msgChan <- fmt.Sprintf("%v", err)
			}
		}()

		f(msgChan, url)
	}
}

func main() {
	urls := os.Args[1:]

	if len(urls) == 0 {
		return
	}

	msgChan := make(chan string)
	mainStart := time.Now()

	for _, url := range urls {
		go protect(msgChan)(
			func(msgChan chan string, url string) {
				requestStart := time.Now()

				resp, err := http.Get(url)
				check(err)

				defer func(Body io.ReadCloser) {
					check(Body.Close())
				}(resp.Body)

				body, err := io.ReadAll(resp.Body)
				check(err)

				requestStop := time.Since(requestStart)
				msgChan <- fmt.Sprintf("%.2fs\t%d\t%s", requestStop.Seconds(), len(body), url)
			},
			url,
		)
	}

	for range urls {
		msg := <-msgChan
		fmt.Println(msg)
	}

	mainStop := time.Since(mainStart)
	fmt.Printf("%.2fs elapsed\n", mainStop.Seconds())
}

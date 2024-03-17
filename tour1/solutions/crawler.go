//go:build exclude

package solutions

import (
	"errors"
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type VisitorState struct {
	state map[string]error
	sync.Mutex
}

func (s *VisitorState) Init() *VisitorState {
	s.state = make(map[string]error)
	return s
}

var (
	visited = (&VisitorState{}).Init()
	loading = errors.New("load in progress")
)

func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		fmt.Printf("<- %v: done\n", url)
		return
	}

	visited.Lock()

	if _, ok := visited.state[url]; ok {
		visited.Unlock()
		fmt.Printf("<- %v: already fetched\n", url)
		return
	}

	visited.state[url] = loading
	visited.Unlock()

	body, urls, err := fetcher.Fetch(url)

	visited.Lock()
	visited.state[url] = err
	visited.Unlock()

	if err != nil {
		fmt.Printf("<- %v: error (%v)\n", url, err)
		return
	}

	fmt.Printf("<-> %v: found (%v)\n", url, body)
	done := make(chan string, len(urls))

	for i, u := range urls {
		fmt.Printf("-> %v: crawling %v/%v (%v)\n", url, i+1, len(urls), u)
		go func(url string) {
			Crawl(url, depth-1, fetcher)
			done <- url
		}(u)
	}

	for i := range urls {
		fmt.Printf("<- %v: crawled %v/%v (%v)\n", url, i+1, len(urls), <-done)
	}

	fmt.Printf("<- %v: done\n", url)
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

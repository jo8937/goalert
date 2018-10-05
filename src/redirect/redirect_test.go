package redirect

import (
	"fmt"
	"sync"
	"testing"
)

// "github.com/patrickmn/go-cache"

var (
	wg sync.WaitGroup
)

func TestReadUrl(t *testing.T) {
	url, err := ReadUrl("aaa")
	if err != nil {
		panic(err)
	}
	fmt.Println(url)

	url, err = ReadUrl("bbb")
	if err != nil {
		panic(err)
	}
	fmt.Println(url)
}

func TestServerPath(t *testing.T) {
	wg.Add(1)
	StartRedirectServer()
	wg.Wait()
}

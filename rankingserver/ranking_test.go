package main

import (
	"fmt"
	"sync"
	"testing"
)

// "github.com/patrickmn/go-cache"

var (
	wg sync.WaitGroup
)

func TestServerPath(t *testing.T) {
	wg.Add(1)
	//StartRedirectServer()
	wg.Wait()
}

func TestFunc(t *testing.T) {
	WriteRanking([]byte(`{"tm":10}`))
	WriteRanking([]byte(`{"tm":5}`))
	WriteRanking([]byte(`{"tm":99}`))
	WriteRanking([]byte(`{"tm":2}`))
	WriteRanking([]byte(`{"tm":45}`))

	js, err := GetRankingJson()

	fmt.Println(err)
	fmt.Println(js)
	fmt.Println("ok")
}

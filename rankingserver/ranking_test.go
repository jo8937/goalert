package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
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

func TestDB(t *testing.T) {
	var err error
	ds := new(DataSource)
	err = ds.Connect()
	if err != nil {
		panic(err)
	}
	defer ds.conn.Close()
	r := Ranking{10, time.Now(), []byte("test")}
	err = ds.InsertRanking(r)
	if err != nil {
		panic(err)
	}

	rankingList, err := ds.ReadRankingList()
	if err != nil {
		panic(err)
	}
	for i, ranking := range rankingList {
		log.Printf("%d %d %s %s \n", i, ranking.sec, ranking.cmt, ranking.regdate)
	}
}

func TestFunc(t *testing.T) {
	WriteRanking([]byte(`{"tm":99181192}`))
	WriteRanking([]byte(`{"tm":99181237}`))
	WriteRanking([]byte(`{"tm":99181186}`))
	WriteRanking([]byte(`{"tm":99181235}`))
	WriteRanking([]byte(`{"tm":99181238}`))
	WriteRanking([]byte(`{"tm":99181194}`))
	WriteRanking([]byte(`{"tm":99181236}`))
	WriteRanking([]byte(`{"tm":99181234}`))

	js, err := GetRankingJson()

	fmt.Println(err)
	fmt.Println(js)
	fmt.Println("ok")
}

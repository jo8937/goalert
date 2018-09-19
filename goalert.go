package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"sync"
	// github.com/patrickmn/go-cache
	// github.com/joho/godotenv
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// main

var (
	debug = flag.Bool("debug", true, "debug")
)

var (
	wg      sync.WaitGroup
	pending chan int
)

//
func reloadKeywords() {

}

// 수동 리로딩 메시지를 받기위한 것
func startManageServer() {
	wg.Done()
}

// filter_exec 에서 끊임없이 stdin 받는 부분
func startReadSTDIN() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "." {
			break
		}
		log.Println(line)
	}
	wg.Done()
}

// 에러 발송
func sendAlert(msg string) {

}

func main() {
	log.Printf("에러로그 리더기 시작")
	flag.Parse()

	wg.Add(1)
	go startManageServer()

	wg.Add(1)
	go startReadSTDIN()

	wg.Wait()
}

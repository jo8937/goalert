package ranking

// https://github.com/valyala/fasthttp/blob/master/examples/helloworldserver/helloworldserver.go

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

var (
	testURL = "http://127.0.0.1:8087"
)

func GetSampleLogJSON() string {
	return `{"tm":99181193}`
}

// func sendRequest(url string, postdata string) {
// 	req, err := http.NewRequest("POST", url, strings.NewReader(postdata))
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Print(resp.StatusCode)

// 	respBody, err := ioutil.ReadAll(resp.Body)
// 	if err == nil {
// 		resText := string(respBody)
// 		log.Printf(resText)
// 	}
// }

func waitForPort(waitFunc func()) {
	for !serverAvailable() {
		time.Sleep(time.Second)
		waitFunc()
	}

}

func InputTest(t *testing.T) {
	resp, err := http.Post(testURL+"/santaserver/regist", "application/json", strings.NewReader(GetSampleLogJSON()))
	if err != nil {
		panic(err)
	}
	t.Log(resp.StatusCode)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("/regist statuscode %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	t.Log(string(data))
}

func ReadTest(t *testing.T) {
	resp, err := http.Get(testURL + "/santaserver/ranking")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("/ranking statuscode %d", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	t.Log(string(data))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func serverAvailable() bool {
	// GET 호출
	resp, err := http.Get(testURL)
	if err != nil {
		log.Print(err)
		return false
	}
	if resp.StatusCode == 200 {
		// 결과 출력
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", string(data))
	}
	return true
}

func Test_WaitServer(t *testing.T) {
	t.Log("start waiting server...")
	waitForPort(func() {
		t.Log("wait a second...")
	})
	t.Log("ok")

}

func Test_Server(t *testing.T) {

	testwg.Add(1)
	go func() {
		StartServer()
		testwg.Done()
	}()

	//time.Sleep(10 * time.Second)

	Test_WaitServer(t)

	InputTest(t)
	ReadTest(t)

	server.Shutdown(nil)

	t.Log("shutdown")

	testwg.Wait()
}

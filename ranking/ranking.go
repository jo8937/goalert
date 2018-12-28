package ranking

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/wangjia184/sortedset"

	_ "github.com/go-sql-driver/mysql"
)

var (
	//rankingSet *SortedSet
	XORKey           = int64(99181225)
	rankingSet       = sortedset.New()
	globalDatasource = new(DataSource)
	server           *http.Server
	wg               sync.WaitGroup
)

// async write db
func WriteRanking(jsonData []byte) (bool, error) {
	var dat map[string]interface{}

	log.Println(string(jsonData))

	if err := json.Unmarshal(jsonData, &dat); err != nil {
		log.Println(err)
		return false, err
	}

	tm, hasKey := dat["tm"]
	if !hasKey {
		log.Printf("error get %s", jsonData)
		return false, errors.New("tm key not found")
	}

	second, castok := tm.(float64)
	if !castok {
		log.Printf("error cast %s", tm)
		return false, errors.New("tm key not int")
	}
	secondInt64 := int64(second)
	secondInt64 = secondInt64 ^ XORKey

	if secondInt64 < 5 || secondInt64 > 3600 {
		log.Printf("error range %d", secondInt64)
		return false, errors.New("tm value in not in range (5~3600)")
	}

	secondString := strconv.FormatInt(secondInt64, 10)

	data := map[string]string{
		"sec":     secondString,
		"regdate": time.Now().Format("2006-01-02 15:04:05"),
	}

	added := rankingSet.AddOrUpdate(secondString, sortedset.SCORE(secondInt64), data)
	return added, nil
}

func SetRankingRegist(body []byte) error {
	log.Printf("%s\n", body)
	return nil
}

func GetRankingJson() (string, error) {
	rankins := rankingSet.GetByRankRange(1, 10, false)
	str, err := json.Marshal(rankins)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

//
func StartServer() {

	globalDatasource.Connect()
	defer globalDatasource.Close()
	uriPrefix := "/santaserver"

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc(uriPrefix+"/regist", func(w http.ResponseWriter, req *http.Request) {
		b, err0 := ioutil.ReadAll(req.Body)
		if err0 != nil {
			w.WriteHeader(500)
			w.Write([]byte("requset body read error"))
		}

		_, err := WriteRanking(b)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("{}"))
		}
		//w.Write([]byte(lastpath))
	})

	http.HandleFunc(uriPrefix+"/ending", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/santa_ending.html", http.StatusSeeOther)
	})

	// response json format [{"Value":{"regdate":"2018-12-26 14:53:05","sec":"99181188"}},{"Value":{"regdate":"2018-12-26 14:53:05","sec":"99181219"}}]
	http.HandleFunc(uriPrefix+"/ranking", func(w http.ResponseWriter, req *http.Request) {
		js, err := GetRankingJson()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(js))
		}
	})

	host := ":8087"
	//server := &http.Server{Addr: ":8087", Handler: s}
	server = &http.Server{Addr: host}
	wg.Add(1)
	go func() {
		// returns ErrServerClosed on graceful close
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop. don't use
			// code with race conditions like these for production. see post
			// comments below on more discussion on how to handle this.
			log.Fatalf("ListenAndServe(): %s", err)
		}
		wg.Done()
	}()
	log.Printf("server started at %s\n", host)
	wg.Wait()

	log.Println("server stopped")
}

func main() {
	flag.Parse()
	StartServer()
}

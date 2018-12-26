package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/wangjia184/sortedset"

	_ "github.com/go-sql-driver/mysql"
)

var (
	//rankingSet *SortedSet
	rankingSet = sortedset.New()
)

func LoadConfig() (string, string, string, int, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	host := os.Getenv("HOST")
	dbname := os.Getenv("DBNAME")

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Port is not number")
	}

	return username, password, host, port, dbname
}

func ReadRankingFromDB(k string) (string, error) {
	username, password, host, port, dbname := LoadConfig()

	// Open database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, dbname))
	if err != nil {
		return "", err // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT URL FROM t_shorten_url WHERE HASH = ?")
	if err != nil {
		return "", err
	}
	defer stmtOut.Close()

	var url string // we "scan" the result in here

	rows, err := stmtOut.Query(k)
	if err != nil {
		return "", err
	}

	if rows.Next() {
		// Query the square-number of 13
		err = rows.Scan(&url) // WHERE number = 13
		if err != nil {
			return "", err
		}
	} else {
		return "", nil
	}

	//fmt.Printf("url : %s", url)
	return url, nil
}

// async write db
func WriteRanking(jsonData []byte) (bool, error) {
	var dat map[string]interface{}

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
	log.Println(rankins)
	str, err := json.Marshal(rankins)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

//
func StartRedirectServer() {
	uri_prefix := "/santaserver"
	http.HandleFunc(uri_prefix+"/regist", func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		err = SetRankingRegist(b)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("{}"))
		}
		//w.Write([]byte(lastpath))
	})

	http.HandleFunc(uri_prefix+"/ending", func(w http.ResponseWriter, req *http.Request) {
		newUrl := "/santa_ending.html"
		http.Redirect(w, req, newUrl, http.StatusSeeOther)
	})

	http.HandleFunc(uri_prefix+"/ranking", func(w http.ResponseWriter, req *http.Request) {

	})

	http.ListenAndServe(":8087", nil)
}

func main() {
	flag.Parse()
}

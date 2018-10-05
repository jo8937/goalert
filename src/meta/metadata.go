package meta

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	// "github.com/go-sql-driver/mysql"
	// "github.com/patrickmn/go-cache"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/patrickmn/go-cache"
)

var (
	c *cache.Cache
)

type Keywords struct {
	Keywords []string `json:"keywords"`
}

func ReadMeta() {
	// Open database connection
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query("SELECT * FROM table")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

}

/*
	메타데이터를 json 에서 다시 읽어서 메모리 cache 에 넣음
*/
func RefreshMeta() {
	c = cache.New(5*time.Minute, 10*time.Minute)
	c.Set("foo", "bar", cache.DefaultExpiration)
	foo, found := c.Get("foo")
	if found {
		fmt.Println(foo)
	}

	byteValues := readJsonConfigFile("keywords.json")
	parseJson(byteValues)
}

/**
파일 전체를 읽어서 byte 배열 리턴
*/
func readJsonConfigFile(filename string) []byte {
	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Panicln(err)
	}

	return byteValue
}

func parseJson(byteValue []byte) {
	var keywords Keywords
	json.Unmarshal(byteValue, &keywords)

	for i := 0; i < len(keywords.Keywords); i++ {
		log.Println("Keyword : " + keywords.Keywords[i])
	}
}

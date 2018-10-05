package redirect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
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

func ReadUrl(k string) (string, error) {
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

func StartRedirectServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		paths := strings.Split(req.URL.Path, "/")
		lastpath := paths[len(paths)-1]

		newUrl, err := ReadUrl(lastpath)

		if err != nil {
			w.Write([]byte(lastpath))
			return
		}
		if newUrl == "" {
			w.Write([]byte(lastpath))
			return
		}

		http.Redirect(w, req, newUrl, http.StatusSeeOther)
	})

	http.ListenAndServe(":8088", nil)
}

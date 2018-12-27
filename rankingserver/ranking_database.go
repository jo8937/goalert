package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

type Ranking struct {
	sec     int
	regdate time.Time
	cmt     []byte
}

// singleton ? connect context
type DataSource struct {
	conn     *sql.DB
	errorLog *log.Logger
}

// db info
type DatabaseConfig struct {
	username string
	password string
	host     string
	port     int
	dbname   string
}

// DB접속. 접속후 conn 변수를 맴버변수로 저장함
func (ds *DataSource) Connect() error {
	dbconf := LoadConfig()
	// Open database connection
	ds.errorLog = log.New(os.Stderr, "### ", log.Ldate|log.Ltime)
	var err error
	ds.conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbconf.username, dbconf.password, dbconf.host, dbconf.port, dbconf.dbname))
	return err
}

// line 을 쿼리로 실행하는 부분
func (ds *DataSource) InsertRanking(r Ranking) error {
	// TODO
	prepared, perr := ds.conn.Prepare("INSERT INTO t_santa_jump_user_ranking(sec, cmt) VALUES(?,?)")
	if perr != nil {
		return perr
	}
	defer prepared.Close()

	insert, qerr := prepared.Query(r.sec, r.cmt)

	if qerr != nil {
		return qerr
	}
	defer insert.Close()
	// be careful deferring Queries if you are using transactions
	return nil
}

func (ds *DataSource) ReadRankingList() ([]Ranking, error) {
	// Prepare statement for reading data
	var rankingList []Ranking

	stmtOut, err := ds.conn.Prepare("SELECT sec, regdate, cmt FROM t_santa_jump_user_ranking ORDER BY sec DESC LIMIT 10")
	if err != nil {
		return nil, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r Ranking
		// Query the square-number of 13
		err = rows.Scan(&r.sec, &r.regdate, &r.cmt) // WHERE number = 13
		if err != nil {
			return nil, err
		}
		rankingList = append(rankingList, r)
	}

	return rankingList, nil
}

// generated from [configure_source_generator.sh]
func LoadConfig() DatabaseConfig {
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

	return DatabaseConfig{username, password, host, port, dbname}
}

package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// 	wg      sync.WaitGroup
	// 	pending chan int
	debug = flag.Bool("debug", false, "debug")
	env   = flag.String("env", "local", "environment (dev/sandbox/production/local)")
)

// singleton ? connect context
type DataSource struct {
	conn       *sql.DB
	successCnt int
	errorCnt   int
	errorLog   *log.Logger
	lineBuffer bytes.Buffer
}

// db info
type DatabaseConfig struct {
	username string
	password string
	host     string
	port     int
	dbname   string
}

func (ds *DataSource) startRead(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		//if len(strings.TrimSpace(line)) == 0 {
		if len(line) == 0 {
			continue
		}
		err := ds.processLine(line)
		if err != nil {
			ds.errorLog.Println(err)
			ds.errorLog.Println(line)
		}
	}
}

// DB접속. 접속후 conn 변수를 맴버변수로 저장함
func (ds *DataSource) Connect() error {
	dbconf := LoadConfig(*env)
	// Open database connection
	ds.errorLog = log.New(os.Stderr, "### ", log.Ldate|log.Ltime)

	ds.successCnt = 0
	ds.errorCnt = 0

	var err error
	ds.conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbconf.username, dbconf.password, dbconf.host, dbconf.port, dbconf.dbname))
	return err
}

// line 을 쿼리로 실행하는 부분
func (ds *DataSource) processLine(line string) error {
	if *debug {
		log.Println(line)
	}

	ds.lineBuffer.WriteString(line)
	// TODO
	if strings.HasSuffix(strings.TrimSpace(line), ";") {
		insert, qerr := ds.conn.Query(ds.lineBuffer.String())
		ds.lineBuffer.Reset()

		if qerr != nil {
			ds.errorCnt++
			return qerr
		}
		ds.successCnt++

		// be careful deferring Queries if you are using transactions
		defer insert.Close()
	}

	return nil
}

func processDatabaseInsert(r io.Reader) error {
	ds := new(DataSource)
	err := ds.Connect()
	if err != nil {
		panic(err)
	}
	defer ds.conn.Close()

	ds.startRead(r)

	log.Printf("Process Finish. Success : %d , Error : %d", ds.successCnt, ds.errorCnt)
	return nil
}

func main() {
	flag.Parse()
	processDatabaseInsert(os.Stdin)
}

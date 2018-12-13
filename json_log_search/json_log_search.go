package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"sync"

	"github.com/google/logger"
	"github.com/json-iterator/go"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// main

var (
	sourcePath    = flag.String("v", "./search_values.txt", "id 목록이 들어있는 라인으로 구분된 값 목록")
	jsonFieldName = flag.String("k", "guid", "JSON 안의 탐색대상 키 필드명")
	targetPath    = flag.String("d", "./data", "JSON 로그들이 들어있는 경로. 이 디랙토리 안의 파일은 모두 탐색 함")
	resultPath    = flag.String("o", "./output/match_data.log", "검색에 걸린 결과를 출력하는 파일")
	concurrent    = flag.Int("n", 8, "동시실행 수 ")
)

var (
	wg               sync.WaitGroup
	fileNameQueue    chan string
	secondaryLogFile *os.File
	secondaryLogger  *log.Logger
	valuelist        map[string]bool

//	fluentdClient    *fluent.Fluent
)

// func initFluentd() {
// 	var conf fluent.Config
// 	conf = fluent.Config{FluentPort: 24224, FluentHost: "localhost", Async: false} // , WriteTimeout: 3
// 	fluentdClient, fluentderr := fluent.New(conf)
// 	if fluentderr != nil {
// 		panic(fluentderr)
// 	}
// }
// func sendFluentd(tag string, data map[string]interface{}) {
// 	flientdErr := fluentdClient.Post(tag, data)
// 	if flientdErr != nil {
// 		// stdout 에 에러로그 찍음
// 		log.Println("fail to send fluentd")
// 		log.Println(flientdErr)
// 		log.Println("----")
// 	}
// }

// 결과를 파일로그로 씀
func createFileLogger() (*log.Logger, *os.File) {
	fl, err := os.OpenFile(*resultPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	return log.New(fl, "", 0), fl
}

func parseJSONLog(bytes []byte) map[string]interface{} {
	var datamap map[string]interface{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(bytes, &datamap); err != nil {
		//log.Printf("JSON Parse ERROR : %s", bytes)
		log.Print(err)
		return nil
	}
	//log.Printf("%+v", datamap)
	//log.Printf("%s", bytes)

	return datamap
}

///////////////////////////////////////////////////
// 로직처리부

func InitReadSourceFile(line string) {
	valuelist[line] = true
}

func CheckJSONLineContainsValues(line string) {
	m := parseJSONLog([]byte(line))
	if m == nil {
		log.Println("line not json : " + line)
	}
	if val, ok := m[*jsonFieldName]; ok {
		strval, strok := val.(string)
		if !strok {
			floatval, floatok := val.(float64)
			if floatok {
				strval = fmt.Sprintf("%.f", floatval)
				strok = true
			}
		}
		if strok {
			// log.Printf("# strval # : %s \n", strval)
			if _, ok2 := valuelist[strval]; ok2 {
				log.Printf("# FOUND # : %s \n", line)
				secondaryLogger.Println(line)
			}
		} else {
			log.Printf("# [%s] field value nil : %s \n", *jsonFieldName, line)
		}

	}
	//log.Printf(line)

}

func ReadEachLine(fileFullPath string, process func(string)) {
	inFile, err := os.Open(fileFullPath)
	defer inFile.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		process(line)
	}
}

// 동시실행 워커
func Worker(filename string) {
	// 큐에서 경쟁적으로 파일을 가져와서 처리
	log.Printf(filename)

	fileFullPath := path.Join(*targetPath, filename)

	ReadEachLine(fileFullPath, CheckJSONLineContainsValues)

	<-fileNameQueue

	wg.Done()
}

// 파일목록을 모두 큐에 넣음
func ListingFilesInDir() {
	files, err := ioutil.ReadDir(*targetPath)
	if err != nil {
		log.Fatal(err)
	}
	fileNameQueue = make(chan string, *concurrent)

	for _, f := range files {
		if !f.IsDir() {
			filename := f.Name()
			fileNameQueue <- filename
			wg.Add(1)
			go Worker(filename)
		}
	}
	close(fileNameQueue)
}

func Execute() {
	valuelist = make(map[string]bool)
	ReadEachLine(*sourcePath, InitReadSourceFile)

	secondaryLogger, secondaryLogFile = createFileLogger()
	defer secondaryLogFile.Close()
	// initFluentd()
	// defer fluentdClient.Close()

	ListingFilesInDir()
	// pending = make(chan int, *concurrent)

	log.Println("Execute finish. Waiting for go rouint end...")
	wg.Wait()

	log.Println("Finish !")
}

///
func main() {
	flag.Parse()
	Execute()
}

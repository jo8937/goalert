package main

import (
	"flag"
	"strings"
	"sync"
	"testing"
)

var (
	wg      sync.WaitGroup
	pending chan int
)

func SetUp() {
	flag.Parse()
	*debug = true
}
func TestEmpty(t *testing.T) {
	SetUp()
	processDatabaseInsert(strings.NewReader(`

`))
}

func TestStdin(t *testing.T) {
	SetUp()
	processDatabaseInsert(strings.NewReader(`

 INSERT INTO aaa ( a,b,c)
     VALUES
		('aa', 'bb', now());

		INSERT INTO aaa ( a,b,c)
		VALUES
		   ('11', '22', now());		
`))

}

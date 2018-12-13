package main

// https://github.com/valyala/fasthttp/blob/master/examples/helloworldserver/helloworldserver.go

import (
	"flag"
	"testing"
)

func TestProcess(t *testing.T) {
	flag.Parse()
	*jsonFieldName = "ddd"
	Execute()
}

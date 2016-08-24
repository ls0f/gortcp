package main

import (
	"flag"
	"gortcp/src"
)

func main() {

	addr := flag.String("addr", "127.0.0.1:33456", "server listen addr")
	flag.Parse()
	c := &gortcp.Client{Addr: *addr}
	c.Connect()
}

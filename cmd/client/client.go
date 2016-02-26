package main

import (
	"flag"

	"github.com/lovedboy/gortcp/src"
)

func main() {

	addr := flag.String("addr", "127.0.0.1:33456", "forward server addr")
	debug := flag.Bool("debug", false, "debug mode(true or false)")
	flag.Parse()
	gortcp.InitLogger(*debug)
	c := &gortcp.Client{Addr: *addr}
	c.Connect()
}
